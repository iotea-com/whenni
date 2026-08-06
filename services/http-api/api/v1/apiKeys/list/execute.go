package apiKeysList

import (
	"fmt"

	sqldb "github.com/ongruent/gruent/db/sqlc"
	gruenthttp "github.com/ongruent/gruent/libs/http"
	"github.com/ongruent/gruent/services/http-api/services/sqlc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func execute(request *gruenthttp.Request[Input]) (*Output, error) {
	request.Span.AddEvent("execute")
	request.Span.SetAttributes(
		attribute.String("request.Input.OrgId", request.Input.OrgId),
		attribute.String("request.Input.SpaceId", request.Input.SpaceId),
		attribute.Int("request.Input.Page", request.Input.Page),
		attribute.Int("request.Input.ResultsPerPage", request.Input.ResultsPerPage),
	)

	dbCtx, dbSpan := otel.Tracer("sqlc").Start(request.Context, "Count API keys")
	var err error
	var apiKeysCount int64
	if request.Input.SpaceId != "" {
		apiKeysCount, err = sqlc.Queries.CountSpaceApiKeys(dbCtx, request.Input.OrgId, &request.Input.SpaceId)
	} else {
		apiKeysCount, err = sqlc.Queries.CountOrganizationApiKeys(dbCtx, request.Input.OrgId)
	}
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error counting API keys in space in the database: %s", err)),
		)
		dbSpan.End()

		return nil, err
	}

	dbSpan.End()

	request.Span.SetAttributes(
		attribute.Int("results.count", int(apiKeysCount)),
	)

	dbCtx, dbSpan = otel.Tracer("sqlc").Start(request.Context, "List API keys")
	var apiKeys []sqldb.AppApiKey
	if request.Input.SpaceId != "" {
		apiKeys, err = sqlc.Queries.ListSpaceApiKeysPaginated(
			dbCtx,
			request.Input.OrgId,
			&request.Input.SpaceId,
			int32(request.Input.ResultsPerPage),
			int32((request.Input.Page-1)*request.Input.ResultsPerPage),
		)
	} else {
		apiKeys, err = sqlc.Queries.ListOrganizationApiKeysPaginated(
			dbCtx,
			request.Input.OrgId,
			int32(request.Input.ResultsPerPage),
			int32((request.Input.Page-1)*request.Input.ResultsPerPage),
		)
	}
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error listing API keys in the database: %s", err)),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	output := Output{
		ApiKeys:        apiKeys,
		Page:           request.Input.Page,
		TotalPages:     (int(apiKeysCount) + request.Input.ResultsPerPage - 1) / request.Input.ResultsPerPage,
		TotalResults:   int(apiKeysCount),
		ResultsPerPage: request.Input.ResultsPerPage,
	}

	return &output, nil
}
