package modelsGet

import (
	"fmt"

	"github.com/jackc/pgx/v5"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	"github.com/iotea-com/iotea/services/http-api/services/sqlc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func execute(request *ioteahttp.Request[Input]) (*Output, error) {
	request.Span.AddEvent("execute")
	request.Span.SetAttributes(
		attribute.String("request.Input.ModelId", request.Input.ModelId),
		attribute.String("request.Input.SpaceId", request.Input.SpaceId),
	)

	dbCtx, dbSpan := otel.Tracer("sqlc").Start(request.Context, "Get model")
	model, err := sqlc.Queries.GetModel(dbCtx, request.Input.ModelId)
	if err != nil {
		if err == pgx.ErrNoRows {
			dbSpan.SetAttributes(
				attribute.String("error.type", "database"),
				attribute.String("error.message", fmt.Sprintf("no model found with ID %s", request.Input.ModelId)),
			)
			dbSpan.End()
			output := Output{
				Model: nil,
			}
			return &output, nil
		}

		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error getting model from the database: %s", err)),
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
