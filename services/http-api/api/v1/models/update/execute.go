package modelsUpdate

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
		attribute.String("request.Input.ModelId", request.Input.ModelId),
		attribute.String("request.Input.SpaceId", request.Input.SpaceId),
	)

	dbCtx, dbSpan := otel.Tracer("sqlc").Start(request.Context, "Update model in space")
	model, err := sqlc.Queries.UpdateModel(
		dbCtx,
		request.Input.ModelId,
		request.Input.Model.Name,
		request.Input.Model.Attributes,
		request.GetActorId(),
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			dbSpan.SetAttributes(
				attribute.String("error.type", "database"),
				attribute.String("error.message", fmt.Sprintf("no model found with ID %s", request.Input.ModelId)),
			)
			dbSpan.End()
			return nil, fiber.NewError(fiber.StatusBadRequest)
		}
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error updating model in the database: %s", err)),
		)
		dbSpan.End()

		return nil, err
	}
	dbSpan.End()

	output := Output{
		Model: &model,
	}

	return &output, nil
}
