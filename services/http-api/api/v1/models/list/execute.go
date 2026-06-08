package modelsList

import (
	"fmt"

	ioteahttp "github.com/iotea-com/iotea/libs/http"
	"github.com/iotea-com/iotea/services/http-api/services/sqlc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func execute(request *ioteahttp.Request[Input]) (*Output, error) {
	request.Span.AddEvent("execute")
	request.Span.SetAttributes(
		attribute.Int("request.Input.Page", request.Input.Page),
		attribute.Int("request.Input.ResultsPerPage", request.Input.ResultsPerPage),
		attribute.String("request.Input.SpaceId", request.Input.SpaceId),
		attribute.String("request.Input.Filter", request.Input.Filter),
	)

	dbCtx, dbSpan := otel.Tracer("sqlc").Start(request.Context, "Count models")
	modelsCount, err := sqlc.Queries.CountModelsFiltered(
		dbCtx,
		request.Input.SpaceId,
		request.Input.Filter,
		request.Input.TagFilter,
	)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error counting models in space in the database: %s", err)),
		)
		dbSpan.End()

		return nil, err
	}

	dbSpan.SetAttributes(
		attribute.Int64("models_count", modelsCount),
	)
	dbSpan.End()

	dbCtx, dbSpan = otel.Tracer("sqlc").Start(request.Context, "List models")
	models, err := sqlc.Queries.ListModelsFiltered(
		dbCtx,
		request.Input.SpaceId,
		request.Input.Filter,
		request.Input.TagFilter,
		int32(request.Input.ResultsPerPage),
		int32((request.Input.Page-1)*request.Input.ResultsPerPage),
	)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error listing models in the database: %s", err)),
		)
		dbSpan.End()

		return nil, err
	}

	dbSpan.End()

	output := &Output{
		Models:         models,
		Page:           request.Input.Page,
		TotalPages:     (int(modelsCount) + request.Input.ResultsPerPage - 1) / request.Input.ResultsPerPage,
		TotalResults:   int(modelsCount),
		ResultsPerPage: request.Input.ResultsPerPage,
	}

	return output, nil
}
