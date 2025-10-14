package permissionsAdd

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
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
		attribute.String("request.Input.SpaceId", request.Input.SpaceId),
		attribute.StringSlice("request.Input.Permissions", request.Input.Permissions),
	)

	// store in the database
	permissionSetId, err := id.Generator.NewPermissionSetId()
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "id_generator"),
			attribute.String("error.message", fmt.Sprintf("error generating a new Permission Set ID: %s", err)),
		)
		request.Span.End()
		return nil, err
	}
	var name *string
	if request.Input.Name != "" {
		name = &request.Input.Name
	}

	var spaceId *string
	if request.Input.SpaceId != "" {
		spaceId = &request.Input.SpaceId
	}

	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Insert permission set")
	permissionSet, err := prisma.Client.PermissionSet.CreateOne(
		db.PermissionSet.ID.Set(*permissionSetId),
		db.PermissionSet.Name.Set(*name),
		db.PermissionSet.CreatedBy.Set(request.GetActorId()),
		db.PermissionSet.UpdatedBy.Set(request.GetActorId()),
		db.PermissionSet.Organization.Link(
			db.Organization.ID.Equals(request.Input.OrgId),
		),
		db.PermissionSet.Space.Link(
			db.Space.ID.EqualsIfPresent(spaceId),
		),
		db.PermissionSet.Permissions.Set(request.Input.Permissions),
	).Exec(dbCtx)

	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error inserting permission set into the database: %s", err)),
		)
		dbSpan.End()
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	dbSpan.End()

	output := Output{
		PermissionSet: permissionSet,
	}

	return &output, nil
}
