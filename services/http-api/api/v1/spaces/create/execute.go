package spacesCreate

import (
	"fmt"

	ioteahttp "github.com/iotea-com/iotea/libs/http"
	ioteapermissions "github.com/iotea-com/iotea/libs/http/permissions"
	"github.com/iotea-com/iotea/libs/id"
	"github.com/iotea-com/iotea/prisma/db"
	"github.com/iotea-com/iotea/services/http-api/services/prisma"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func execute(request *ioteahttp.Request[Input]) (*Output, error) {
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

	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Insert space")
	space, err := prisma.Client.Space.CreateOne(
		db.Space.ID.Set(*spaceId),
		db.Space.Name.Set(request.Input.Name),
		db.Space.CreatedBy.Set(request.GetActorId()),
		db.Space.UpdatedBy.Set(request.GetActorId()),
		db.Space.Organization.Link(
			db.Organization.ID.Equals(request.Input.OrgId),
		),
	).Exec(dbCtx)
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

	dbCtx, dbSpan = otel.Tracer("prisma").Start(request.Context, "Insert default permission set")
	_, err = prisma.Client.PermissionSet.CreateOne(
		db.PermissionSet.ID.Set(*defaultPermissionSetId),
		db.PermissionSet.Name.Set("Default"),
		db.PermissionSet.CreatedBy.Set(request.GetActorId()),
		db.PermissionSet.UpdatedBy.Set(request.GetActorId()),
		db.PermissionSet.Organization.Link(
			db.Organization.ID.Equals(space.OrganizationID),
		),
		db.PermissionSet.Space.Link(
			db.Space.ID.Equals(space.ID),
		),
		db.PermissionSet.Permissions.Set(ioteapermissions.DefaultMemberSpacePermissions),
	).Exec(dbCtx)
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
		Space: space,
	}

	return &output, nil
}
