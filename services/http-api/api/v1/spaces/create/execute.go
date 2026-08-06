package spacesCreate

import (
	"fmt"

	sqldb "github.com/ongruent/gruent/db/sqlc"
	gruenthttp "github.com/ongruent/gruent/libs/http"
	gruentpermissions "github.com/ongruent/gruent/libs/http/permissions"
	"github.com/ongruent/gruent/libs/id"
	"github.com/ongruent/gruent/services/http-api/services/sqlc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func execute(request *gruenthttp.Request[Input]) (*Output, error) {
	request.Span.AddEvent("execute")
	request.Span.SetAttributes(
		attribute.String("request.Input.Name", request.Input.Name),
		attribute.String("request.Input.OrgId", request.Input.OrgId),
	)

	// Create space
	spaceId, err := id.Generator.NewSpaceId()
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "id_generator"),
			attribute.String("error.message", fmt.Sprintf("error generating a new Space ID: %s", err)),
		)
		return nil, err
	}

	dbCtx, dbSpan := otel.Tracer("sqlc").Start(request.Context, "Insert space")
	space, err := sqlc.Queries.CreateSpace(
		dbCtx,
		*spaceId,
		request.Input.OrgId,
		request.Input.Name,
		request.GetActorId(),
		request.GetActorId(),
	)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error inserting space into the database: %s", err)),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	// Create default permission set
	defaultPermissionSetId, err := id.Generator.NewPermissionSetId()
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "id_generator"),
			attribute.String("error.message", fmt.Sprintf("error generating a new Default Permission Set ID: %s", err)),
		)
		return nil, err
	}

	dbCtx, dbSpan = otel.Tracer("sqlc").Start(request.Context, "Insert default permission set")
	_, err = sqlc.Queries.CreatePermissionSet(dbCtx, sqldb.CreatePermissionSetParams{
		ID:             *defaultPermissionSetId,
		OrganizationID: space.OrganizationID,
		SpaceID:        &space.ID,
		Name:           "Default",
		Permissions:    gruentpermissions.DefaultMemberSpacePermissions,
		CreatedBy:      request.GetActorId(),
		UpdatedBy:      request.GetActorId(),
	})
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error inserting default permission set into the database: %s", err)),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	output := Output{
		Space: &space,
	}

	return &output, nil
}
