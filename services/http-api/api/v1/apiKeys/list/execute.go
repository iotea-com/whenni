package apiKeysList

import (
	"fmt"
	"strconv"

	ioteahttp "github.com/iotea-com/iotea/libs/http"
	"github.com/iotea-com/iotea/prisma/db"
	"github.com/iotea-com/iotea/services/http-api/services/prisma"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func execute(request *ioteahttp.Request[Input]) (*Output, error) {
	request.Span.AddEvent("execute")
	request.Span.SetAttributes(
		attribute.String("request.Input.OrgId", request.Input.OrgId),
		attribute.String("request.Input.SpaceId", request.Input.SpaceId),
		attribute.Int("request.Input.Page", request.Input.Page),
		attribute.Int("request.Input.ResultsPerPage", request.Input.ResultsPerPage),
	)

	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Count API keys")
	var apiKeysCountResponse []struct {
		Count db.RawString `json:"api_keys_count"`
	}
	countQuery := `SELECT count(*) as api_keys_count FROM app."apiKeys" WHERE "apiKeys"."organizationId" = $1 AND "apiKeys"."spaceId" IS NULL`
	if request.Input.SpaceId != "" {
		countQuery = `SELECT count(*) as api_keys_count FROM app."apiKeys" WHERE "apiKeys"."organizationId" = $1 AND "apiKeys"."spaceId" = $2`
	}
	err := prisma.Client.Prisma.QueryRaw(countQuery, request.Input.OrgId, request.Input.SpaceId).Exec(dbCtx, &apiKeysCountResponse)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error counting API keys in space in the database: %s", err)),
		)
		dbSpan.End()

		return nil, err
	}

	dbSpan.End()

	apiKeysCount, err := strconv.ParseInt(string(apiKeysCountResponse[0].Count), 10, 16)
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "parse_int"),
			attribute.String("error.message", fmt.Sprintf("could not parse count response to int - %#v: %s", apiKeysCountResponse, err)),
		)

		return nil, err
	}

	request.Span.SetAttributes(
		attribute.Int("results.count", int(apiKeysCount)),
	)
	dbSpan.End()

	dbCtx, dbSpan = otel.Tracer("prisma").Start(request.Context, "List API keys")
	params := []db.APIKeyWhereParam{
		db.APIKey.OrganizationID.Equals(request.Input.OrgId),
	}
	if request.Input.SpaceId != "" {
		params = append(params, db.APIKey.SpaceID.Equals(request.Input.SpaceId))
	} else {
		params = append(params, db.APIKey.SpaceID.IsNull())
	}
	apiKeys, err := prisma.Client.APIKey.FindMany(
		params...,
	).With(
		db.APIKey.OrganizationPermissionSet.Fetch(),
		db.APIKey.SpacePermissionSet.Fetch(),
	).OrderBy(db.APIKey.CreatedAt.Order(db.DESC)).
		Take(request.Input.ResultsPerPage).
		Skip((request.Input.Page - 1) * request.Input.ResultsPerPage).
		Exec(dbCtx)
	if err != nil {
		if err.Error() == "ErrNotFound" {
			dbSpan.SetAttributes(
				attribute.String("error.type", "database"),
				attribute.String("error.message", fmt.Sprintf("no API keys found in space with ID %s", request.Input.SpaceId)),
			)
			dbSpan.End()

			output := Output{
				ApiKeys: []db.APIKeyModel{},
			}

			return &output, nil
		}

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
