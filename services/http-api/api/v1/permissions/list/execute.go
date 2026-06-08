package permissionsList

import (
	"fmt"

	sqldb "github.com/iotea-com/iotea/db/sqlc"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	"github.com/iotea-com/iotea/services/http-api/services/sqlc"
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

	dbCtx, dbSpan := otel.Tracer("sqlc").Start(request.Context, "Count permission sets")
	var permissionSetsCount int64
	var err error
	if request.Input.SpaceId != "" {
		permissionSetsCount, err = sqlc.Queries.CountSpacePermissionSets(dbCtx, &request.Input.SpaceId)
	} else {
		permissionSetsCount, err = sqlc.Queries.CountOrganizationPermissionSets(dbCtx, request.Input.OrgId)
	}
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error counting permissions in space in the database: %s", err)),
		)
		dbSpan.End()

		return nil, err
	}

	dbSpan.SetAttributes(attribute.Int("permission_sets_count", int(permissionSetsCount)))
	dbSpan.End()

	dbCtx, dbSpan = otel.Tracer("sqlc").Start(request.Context, "List permission sets in scope")
	offset := (request.Input.Page - 1) * request.Input.ResultsPerPage
	var permissionSets []sqldb.AppPermission
	if request.Input.SpaceId != "" {
		permissionSets, err = sqlc.Queries.ListSpacePermissionSetsPaginated(
			dbCtx,
			&request.Input.SpaceId,
			int32(request.Input.ResultsPerPage),
			int32(offset),
		)
	} else {
		permissionSets, err = sqlc.Queries.ListOrganizationPermissionSetsPaginated(
			dbCtx,
			request.Input.OrgId,
			int32(request.Input.ResultsPerPage),
			int32(offset),
		)
	}
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error listing permission sets in the database: %s", err)),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()
	if permissionSets == nil {
		permissionSets = []sqldb.AppPermission{}
	}

	output := Output{
		PermissionSets: permissionSets,
		Page:           request.Input.Page,
		TotalPages:     (int(permissionSetsCount) + request.Input.ResultsPerPage - 1) / request.Input.ResultsPerPage,
		TotalResults:   int(permissionSetsCount),
		ResultsPerPage: request.Input.ResultsPerPage,
	}

	return &output, nil
}
