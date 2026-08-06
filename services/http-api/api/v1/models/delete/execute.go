package modelsDelete

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	gruenthttp "github.com/ongruent/gruent/libs/http"
	gruentchannels "github.com/ongruent/gruent/libs/legacy/engine/channels"
	"github.com/ongruent/gruent/services/http-api/services/sqlc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func execute(request *gruenthttp.Request[Input]) (*Output, error) {
	request.Span.AddEvent("execute")
	request.Span.SetAttributes(
		attribute.String("request.Input.ModelId", request.Input.ModelId),
		attribute.String("request.Input.SpaceId", request.Input.SpaceId),
	)

	// Check if model is used in any channel
	dbCtx, dbSpan := otel.Tracer("sqlc").Start(request.Context, "Check if model is used in channel")
	channels, err := sqlc.Queries.ListChannels(dbCtx, request.Input.SpaceId)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error checking if model is used in a channel: %s", err)),
		)
		dbSpan.End()

		return nil, err
	}

	dbSpan.End()

	channelsWithModel := []string{}
	for _, channel := range channels {
		var channelConfig gruentchannels.Channel
		err = json.Unmarshal(channel.Config, &channelConfig)
		if err != nil {
			request.Span.SetAttributes(
				attribute.String("error.type", "parse_json"),
				attribute.String("error.message", fmt.Sprintf("could not read channel config: %s", err.Error())),
			)

			errorResponse := gruenthttp.NewErrorResponse([]any{
				fmt.Sprintf("could not read channel config: %s", err.Error()),
			})
			errorResponseJson, _ := errorResponse.MarshalJson()
			return nil, fiber.NewError(fiber.StatusBadRequest, string(errorResponseJson))
		}

		for _, node := range channelConfig.Nodes {
			for _, modelDependency := range node.Metadata.Dependencies.Models {

				if modelDependency.ModelId == request.Input.ModelId {
					channelsWithModel = append(channelsWithModel, channel.ID)
				}
			}
		}
	}

	if len(channelsWithModel) > 0 {
		errMessage := fmt.Sprintf("model is used in %d channel(s)", len(channelsWithModel))
		request.Span.SetAttributes(
			attribute.String("error.type", "context_validation"),
			attribute.String("error.message", errMessage),
		)

		errorResponse := gruenthttp.NewErrorResponse([]any{
			errMessage,
		})
		errorResponseJson, _ := errorResponse.MarshalJson()
		return nil, fiber.NewError(fiber.StatusConflict, string(errorResponseJson))
	}

	// Delete thing
	dbCtx, dbSpan = otel.Tracer("sqlc").Start(request.Context, "Delete model")
	_, err = sqlc.Queries.GetModel(dbCtx, request.Input.ModelId)
	if err != nil {
		if err == pgx.ErrNoRows {
			dbSpan.SetAttributes(
				attribute.String("error.type", "database"),
				attribute.String("error.message", fmt.Sprintf("no model found with ID %s in space with ID %s", request.Input.ModelId, request.Input.SpaceId)),
			)
			dbSpan.End()
			return nil, fiber.NewError(fiber.StatusBadRequest)
		}
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error checking model in the database: %s", err)),
		)
		dbSpan.End()
		return nil, err
	}

	err = sqlc.Queries.DeleteModel(dbCtx, request.Input.ModelId)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error deleting model from the database: %s", err)),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	return nil, nil
}
