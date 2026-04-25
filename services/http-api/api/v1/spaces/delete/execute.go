package spacesDelete

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	"github.com/iotea-com/iotea/services/http-api/services/sqlc"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func execute(request *ioteahttp.Request[Input]) (*Output, error) {
	request.Span.AddEvent("execute")
	request.Span.SetAttributes(
		attribute.String("request.Input.SpaceId", request.Input.SpaceId),
	)

	spaceId := request.Input.SpaceId

	dbCtx, dbSpan := otel.Tracer("sqlc").Start(request.Context, "Get space")
	_, err := sqlc.Queries.GetSpace(dbCtx, spaceId)
	if err != nil {
		if err == pgx.ErrNoRows {
			dbSpan.SetAttributes(
				attribute.String("error.type", "database"),
				attribute.String("error.message", fmt.Sprintf("no space found with ID %s", spaceId)),
			)
			dbSpan.End()
			return nil, fiber.NewError(fiber.StatusBadRequest)
		}
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error getting space from the database: %s", err)),
		)
		dbSpan.End()
		return nil, err
	}
	dbSpan.End()

	// Get published channels for this space from the database
	dbCtx, dbSpan = otel.Tracer("sqlc").Start(request.Context, "Get published channels")
	publishedChannels, err := sqlc.Queries.ListPublishedChannels(dbCtx, spaceId)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error getting published channels from the database: %s", err)),
		)
		dbSpan.End()
		return nil, err
	}
	dbSpan.End()

	// Do not delete space if there are published channels
	if len(publishedChannels) > 0 {
		request.Span.SetAttributes(
			attribute.String("error.type", "context_validation"),
			attribute.String("error.message", "delete space failed since there are published channels linked to the space"),
		)
		response := ioteahttp.NewErrorResponse([]any{"unpublish or delete all channels in this space before deleting the space"})
		responseJson, _ := response.MarshalJson()
		return nil, fiber.NewError(fiber.StatusConflict, string(responseJson))
	}

	// Delete models
	dbCtx, dbSpan = otel.Tracer("sqlc").Start(request.Context, "Delete models")
	err = sqlc.Queries.DeleteModelsBySpace(dbCtx, spaceId)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error deleting models from the database: %s", err)),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	// Delete things
	dbCtx, dbSpan = otel.Tracer("sqlc").Start(request.Context, "Delete things")
	err = sqlc.Queries.DeleteThingsBySpace(dbCtx, spaceId)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error deleting things from the database: %s", err)),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	// Delete certificates
	dbCtx, dbSpan = otel.Tracer("sqlc").Start(request.Context, "Delete certificates")
	err = sqlc.Queries.DeleteCertificatesBySpace(dbCtx, spaceId)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error deleting certificates from the database: %s", err)),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	// Delete api keys
	dbCtx, dbSpan = otel.Tracer("sqlc").Start(request.Context, "Delete API keys")
	err = sqlc.Queries.DeleteApiKeysBySpace(dbCtx, &spaceId)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error deleting API keys from the database: %s", err)),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	// Delete permission sets
	dbCtx, dbSpan = otel.Tracer("sqlc").Start(request.Context, "Delete permission sets")
	err = sqlc.Queries.DeletePermissionSetsBySpace(dbCtx, &spaceId)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error deleting permission sets from the database: %s", err)),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	// Delete channels
	dbCtx, dbSpan = otel.Tracer("sqlc").Start(request.Context, "Delete channels")
	err = sqlc.Queries.DeleteChannelsBySpace(dbCtx, request.Input.SpaceId)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error deleting channels from the database: %s", err)),
		)
		dbSpan.End()

		return nil, err
	}

	dbSpan.End()

	// Delete space
	dbCtx, dbSpan = otel.Tracer("sqlc").Start(request.Context, "Delete space")
	err = sqlc.Queries.DeleteSpace(dbCtx, spaceId)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error deleting space from the database: %s", err)),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	output := Output{}

	return &output, nil
}
