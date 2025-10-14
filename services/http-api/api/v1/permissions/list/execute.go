package permissionsList

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
		attribute.String("request.Input.ChannelId", request.Input.OrgId),
		attribute.String("request.Input.SpaceId", request.Input.SpaceId),
		attribute.Int("request.Input.Page", request.Input.Page),
		attribute.Int("request.Input.ResultsPerPage", request.Input.ResultsPerPage),
	)

	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Count permission sets")
	var permissionSetsCountResponse []struct {
		Count db.RawString `json:"permission_sets_count"`
	}
	countQuery := `SELECT count(*) as permission_sets_count FROM app.permissions WHERE permissions."organizationId" = $1 AND permissions."spaceId" IS NULL`
	if request.Input.SpaceId != "" {
		countQuery = `SELECT count(*) as permission_sets_count FROM app.permissions WHERE permissions."spaceId" = $1`
	}
	scopeId := request.Input.OrgId
	if request.Input.SpaceId != "" {
		scopeId = request.Input.SpaceId
	}
	err := prisma.Client.Prisma.QueryRaw(countQuery, scopeId).Exec(dbCtx, &permissionSetsCountResponse)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error counting permissions in space in the database: %s", err)),
		)
		dbSpan.End()

		return nil, err
	}

	permissionSetsCount, err := strconv.ParseInt(string(permissionSetsCountResponse[0].Count), 10, 16)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("could not parse count response to int - %#v: %s", permissionSetsCountResponse, err)),
		)
		dbSpan.End()

		return nil, err
	}

	dbSpan.SetAttributes(attribute.Int("permission_sets_count", int(permissionSetsCount)))
	dbSpan.End()

	dbCtx, dbSpan = otel.Tracer("prisma").Start(request.Context, "List permission sets in the space")

	permissionSetQueryParams := []db.PermissionSetWhereParam{}
	if request.Input.SpaceId != "" {
		permissionSetQueryParams = append(permissionSetQueryParams, db.PermissionSet.SpaceID.Equals(request.Input.SpaceId))
	} else {
		permissionSetQueryParams = append(permissionSetQueryParams, db.PermissionSet.OrganizationID.Equals(request.Input.OrgId))
		permissionSetQueryParams = append(permissionSetQueryParams, db.PermissionSet.SpaceID.IsNull())
	}

	permissionSets, err := prisma.Client.PermissionSet.FindMany(
		permissionSetQueryParams...,
	).OrderBy(db.PermissionSet.UpdatedAt.Order(db.DESC)).
		Take(request.Input.ResultsPerPage).
		Skip((request.Input.Page - 1) * request.Input.ResultsPerPage).
		Exec(dbCtx)
	if err != nil {
		if err.Error() == "ErrNotFound" {
			dbSpan.SetAttributes(
				attribute.String("error.type", "database"),
				attribute.String("error.message", fmt.Sprintf("no permission sets found in space with ID %s", request.Input.SpaceId)),
			)
			dbSpan.End()

			output := Output{
				PermissionSets: []db.PermissionSetModel{},
			}

			return &output, nil
		}

		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error listing permission sets in the database: %s", err)),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	output := Output{
		PermissionSets: permissionSets,
		Page:           request.Input.Page,
		TotalPages:     (int(permissionSetsCount) + request.Input.ResultsPerPage - 1) / request.Input.ResultsPerPage,
		TotalResults:   int(permissionSetsCount),
		ResultsPerPage: request.Input.ResultsPerPage,
	}

	return &output, nil
}
