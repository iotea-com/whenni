package channelsDelete

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

	// get channel
	dbCtx, dbSpan := otel.Tracer("sqlc").Start(request.Context, "Get channel")
	channel, err := sqlc.Queries.GetChannel(dbCtx, request.Input.ChannelId)
	if err != nil {
		if err == pgx.ErrNoRows {
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
	if channel.PublishedAt.Valid || channel.PublishedBy != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "context_validation"),
			attribute.String("error.message", fmt.Sprintf("no channel found with ID %s", request.Input.ChannelId)),
		)
		errorResponse := gruenthttp.NewErrorResponse([]any{
			"could not delete a published channel",
		})
		errorResponseJson, _ := errorResponse.MarshalJson()
		return nil, fiber.NewError(fiber.StatusConflict, string(errorResponseJson))
	}

	// delete channel
	dbCtx, dbSpan = otel.Tracer("sqlc").Start(request.Context, "Delete channel")
	err = sqlc.Queries.DeleteChannel(dbCtx, request.Input.ChannelId)
	if err != nil {
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
