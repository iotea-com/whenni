package channelsGet

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	gruenthttp "github.com/ongruent/gruent/libs/http"
	"github.com/ongruent/gruent/services/http-api/services/sqlc"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func execute(request *gruenthttp.Request[Input]) (*Output, error) {
	request.Span.AddEvent("execute")
	request.Span.SetAttributes(
		attribute.String("request.Input.ChannelId", request.Input.ChannelId),
		attribute.String("request.Input.SpaceId", request.Input.SpaceId),
	)

	dbCtx, dbSpan := otel.Tracer("sqlc").Start(request.Context, "Get channel")
	channel, err := sqlc.Queries.GetChannel(dbCtx, request.Input.ChannelId)
	if err != nil {
		if err == pgx.ErrNoRows {
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
		Channel: &channel,
	}

	return &output, nil
}
