package modelsDelete

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	ioteachannels "github.com/iotea-com/iotea/libs/engine/channels"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	"github.com/iotea-com/iotea/prisma/db"
	"github.com/iotea-com/iotea/services/http-api/services/prisma"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func execute(request *ioteahttp.Request[Input]) (*Output, error) {
	request.Span.AddEvent("execute")
	request.Span.SetAttributes(
		attribute.String("request.Input.ModelId", request.Input.ModelId),
		attribute.String("request.Input.SpaceId", request.Input.SpaceId),
	)

	// Check if model is used in any channel
	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Check if model is used in channel")
	channels, err := prisma.Client.Channel.FindMany(
		db.Channel.SpaceID.Equals(request.Input.SpaceId),
	).Select(
		db.Channel.Config.Field(),
	).Exec(dbCtx)
	if err != nil && err.Error() != "ErrNotFound" {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error checking if model is used in a channel: %s", err)),
		)
		dbSpan.End()

		return nil, err
	}

	dbSpan.End()

	channelsWithModel := []db.ChannelModel{}
	for _, channel := range channels {
		var channelConfig ioteachannels.Channel
		err = json.Unmarshal(channel.Config, &channelConfig)
		if err != nil {
			request.Span.SetAttributes(
				attribute.String("error.type", "parse_json"),
				attribute.String("error.message", fmt.Sprintf("could not read channel config: %s", err.Error())),
			)

			errorResponse := ioteahttp.NewErrorResponse([]any{
				fmt.Sprintf("could not read channel config: %s", err.Error()),
			})
			errorResponseJson, _ := errorResponse.MarshalJson()
			return nil, fiber.NewError(fiber.StatusBadRequest, string(errorResponseJson))
		}

		for _, node := range channelConfig.Nodes {
			for _, modelDependency := range node.Metadata.Dependencies.Models {

				if modelDependency.ModelId == request.Input.ModelId {
					channelsWithModel = append(channelsWithModel, channel)
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

		errorResponse := ioteahttp.NewErrorResponse([]any{
			errMessage,
		})
		errorResponseJson, _ := errorResponse.MarshalJson()
		return nil, fiber.NewError(fiber.StatusConflict, string(errorResponseJson))
	}

	// Delete thing
	dbCtx, dbSpan = otel.Tracer("prisma").Start(request.Context, "Delete model")
	_, err = prisma.Client.Model.FindUnique(
		db.Model.ID.Equals(request.Input.ModelId),
	).Delete().Exec(dbCtx)
	if err != nil {
		if err.Error() == "ErrNotFound" {
			dbSpan.SetAttributes(
				attribute.String("error.type", "database"),
				attribute.String("error.message", fmt.Sprintf("no model found with ID %s in space with ID %s", request.Input.ModelId, request.Input.SpaceId)),
			)
			dbSpan.End()
			return nil, fiber.NewError(fiber.StatusBadRequest)
		}

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
