package thingsDelete

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	gruenthttp "github.com/ongruent/gruent/libs/http"
	gruentchannels "github.com/ongruent/gruent/libs/legacy/engine/channels"
	"github.com/ongruent/gruent/services/http-api/services/sqlc"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func execute(request *gruenthttp.Request[Input]) (*Output, error) {
	request.Span.AddEvent("execute")
	request.Span.SetAttributes(
		attribute.String("request.Input.ThingId", request.Input.ThingId),
		attribute.String("request.Input.SpaceId", request.Input.SpaceId),
	)

	// Check if thing is used in any channel
	dbCtx, dbSpan := otel.Tracer("sqlc").Start(request.Context, "Check if thing is used in channel")
	channels, err := sqlc.Queries.ListChannels(dbCtx, request.Input.SpaceId)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error checking if thing is used in a channel: %s", err)),
		)
		dbSpan.End()

		return nil, err
	}

	dbSpan.End()

	channelsWithThing := []string{}
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
			for _, thingDependency := range node.Metadata.Dependencies.Things {

				if thingDependency.ThingId == request.Input.ThingId {
					channelsWithThing = append(channelsWithThing, channel.ID)
				}
			}
		}
	}

	if len(channelsWithThing) > 0 {
		errMessage := fmt.Sprintf("thing is used in %d channel(s)", len(channelsWithThing))
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

	// Check if thing is referenced in any other things
	dbCtx, dbSpan = otel.Tracer("sqlc").Start(request.Context, "Check if thing is referenced in another thing")
	things, err := sqlc.Queries.ListThings(dbCtx, request.Input.SpaceId)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error checking if thing is used in another thing: %s", err)),
		)
		dbSpan.End()

		return nil, err
	}

	dbSpan.End()

	thingsWithThing := []string{}
	for _, thing := range things {
		var thingAttributes map[string]any
		err = json.Unmarshal(thing.Attributes, &thingAttributes)
		if err != nil {
			request.Span.SetAttributes(
				attribute.String("error.type", "parse_json"),
				attribute.String("error.message", fmt.Sprintf("could not read thing attributes: %s", err.Error())),
			)

			errorResponse := gruenthttp.NewErrorResponse([]any{
				fmt.Sprintf("could not read thing attributes: %s", err.Error()),
			})
			errorResponseJson, _ := errorResponse.MarshalJson()
			return nil, fiber.NewError(fiber.StatusBadRequest, string(errorResponseJson))
		}

		for _, value := range thingAttributes {
			if value == request.Input.ThingId {
				thingsWithThing = append(thingsWithThing, thing.ID)
			}
		}
	}

	if len(thingsWithThing) > 0 {
		errMessage := fmt.Sprintf("thing is referenced in %d other thing(s)", len(thingsWithThing))
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
	dbCtx, dbSpan = otel.Tracer("sqlc").Start(request.Context, "Delete thing")
	_, err = sqlc.Queries.GetThing(dbCtx, request.Input.ThingId)
	if err != nil {
		if err == pgx.ErrNoRows {
			dbSpan.SetAttributes(
				attribute.String("error.type", "database"),
				attribute.String("error.message", fmt.Sprintf("no thing found with ID %s", request.Input.ThingId)),
			)
			dbSpan.End()

			return nil, fiber.NewError(fiber.StatusBadRequest)
		}

		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error deleting thing from the database: %s", err)),
		)
		dbSpan.End()

		return nil, err
	}

	err = sqlc.Queries.DeleteThing(dbCtx, request.Input.ThingId)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error deleting thing from the database: %s", err)),
		)
		dbSpan.End()

		return nil, err
	}

	dbSpan.End()

	return nil, nil
}
