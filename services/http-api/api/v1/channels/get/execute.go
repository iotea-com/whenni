package channelsGet

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

	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Get channel")
	channel, err := prisma.Client.Channel.FindUnique(
		db.Channel.ID.Equals(request.Input.ChannelId),
	).Exec(dbCtx)
	if err != nil {
		if err.Error() == "ErrNotFound" {
			request.Span.SetAttributes(
				attribute.String("error.type", "http_request"),
				attribute.String("error.message", fmt.Sprintf("no channel found with ID %s", request.Input.ChannelId)),
			)
			dbSpan.End()

			return nil, fiber.NewError(fiber.StatusBadRequest)
		}

		request.Span.SetAttributes(
			attribute.String("error.type", "http_request"),
			attribute.String("error.message", fmt.Sprintf("error getting channel from the database: %s", err)),
		)
		dbSpan.End()

		return nil, err
	}

	dbSpan.End()

	output := Output{
		Channel: channel,
	}

	return &output, nil
}
