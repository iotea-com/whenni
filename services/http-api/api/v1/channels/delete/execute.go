package channelsDelete

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	"github.com/iotea-com/iotea/prisma/db"
	"github.com/iotea-com/iotea/services/http-api/services/prisma"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func execute(request *ioteahttp.Request[Input]) (*Output, error) {
	request.Span.AddEvent("execute")
	request.Span.SetAttributes(
		attribute.String("request.Input.ChannelId", request.Input.ChannelId),
		attribute.String("request.Input.SpaceId", request.Input.SpaceId),
	)

	// get channel
	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Get channel")
	channel, err := prisma.Client.Channel.FindUnique(
		db.Channel.ID.Equals(request.Input.ChannelId),
	).Exec(dbCtx)
	if err != nil {
		if err.Error() == "ErrNotFound" {
			dbSpan.SetAttributes(
				attribute.String("error.type", "database"),
				attribute.String("error.message", fmt.Sprintf("no channel found with ID %s", request.Input.ChannelId)),
			)
			dbSpan.End()
			return nil, fiber.NewError(fiber.StatusBadRequest)
		}

		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("could not get channel from the database: %s", err)),
		)
		dbSpan.End()

		return nil, err
	}

	dbSpan.End()

	// verify channel is not published
	_, hasPublishedAt := channel.PublishedAt()
	_, hasPublishedBy := channel.PublishedBy()
	if hasPublishedAt || hasPublishedBy {
		dbSpan.SetAttributes(
			attribute.String("error.type", "context_validation"),
			attribute.String("error.message", fmt.Sprintf("no channel found with ID %s", request.Input.ChannelId)),
		)
		errorResponse := ioteahttp.NewErrorResponse([]any{
			"could not delete a published channel",
		})
		errorResponseJson, _ := errorResponse.MarshalJson()
		return nil, fiber.NewError(fiber.StatusConflict, string(errorResponseJson))
	}

	// delete channel
	dbCtx, dbSpan = otel.Tracer("prisma").Start(request.Context, "Delete channel")
	_, err = prisma.Client.Channel.FindUnique(
		db.Channel.ID.Equals(request.Input.ChannelId),
	).Delete().Exec(dbCtx)
	if err != nil {
		if err.Error() == "ErrNotFound" {
			dbSpan.SetAttributes(
				attribute.String("error.type", "database"),
				attribute.String("error.message", fmt.Sprintf("no channel found with ID %s", request.Input.ChannelId)),
			)

			return nil, fiber.NewError(fiber.StatusBadRequest)
		}

		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("could not delete channel from the database: %s", err)),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	return nil, nil
}
