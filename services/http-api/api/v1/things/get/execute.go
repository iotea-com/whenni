package thingsGet

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
		attribute.String("request.Input.ThingId", request.Input.ThingId),
		attribute.String("request.Input.SpaceId", request.Input.SpaceId),
	)

	dbCtx, dbSpan := otel.Tracer("sqlc").Start(request.Context, "Get thing")
	thing, err := sqlc.Queries.GetThing(dbCtx, request.Input.ThingId)
	if err != nil {
		if err == pgx.ErrNoRows {
			err = fmt.Errorf("no thing found with ID %s", request.Input.ThingId)
			dbSpan.SetAttributes(
				attribute.String("error.type", "database"),
				attribute.String("error.message", err.Error()),
			)
			dbSpan.End()

			return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error getting thing from the database: %s", err)),
		)
		dbSpan.End()

		return nil, err
	}

	// if thing is internal, do not return it
	if thing.Internal {
		errMessage := "error getting thing from the database: thing is internal"
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", errMessage),
		)
		dbSpan.End()

		return nil, fiber.NewError(fiber.StatusBadRequest, errMessage)
	}

	dbSpan.End()

	output := Output{
		Thing: &thing,
	}

	return &output, nil
}
