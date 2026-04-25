package policiesList

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
		attribute.String("request.Input.SpaceId", request.Input.SpaceId),
		attribute.Int("request.Input.Page", request.Input.Page),
		attribute.Int("request.Input.ResultsPerPage", request.Input.ResultsPerPage),
	)

	dbCtx, dbSpan := otel.Tracer("sqlc").Start(request.Context, "Count policies")
	policiesCount, err := sqlc.Queries.CountPolicies(dbCtx, request.Input.SpaceId)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error counting policies in space in the database: %s", err)),
		)
		dbSpan.End()

		return nil, err
	}

	dbSpan.SetAttributes(
		attribute.Int64("policiesCount", policiesCount),
	)
	dbSpan.End()

	dbCtx, dbSpan = otel.Tracer("sqlc").Start(request.Context, "List policies in the space")
	certificates, err := sqlc.Queries.ListPoliciesPaginated(
		dbCtx,
		request.Input.SpaceId,
		int32(request.Input.ResultsPerPage),
		int32((request.Input.Page-1)*request.Input.ResultsPerPage),
	)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error listing permission sets in the database: %s", err)),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	output := Output{
		Certificates:   certificates,
		Page:           request.Input.Page,
		TotalPages:     (int(policiesCount) + request.Input.ResultsPerPage - 1) / request.Input.ResultsPerPage,
		TotalResults:   int(policiesCount),
		ResultsPerPage: request.Input.ResultsPerPage,
	}

	return &output, nil
}
