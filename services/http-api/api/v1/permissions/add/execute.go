package permissionsAdd

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	sqldb "github.com/iotea-com/iotea/db/sqlc"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	"github.com/iotea-com/iotea/libs/id"
	"github.com/iotea-com/iotea/services/http-api/services/sqlc"
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
	var spaceId *string
	if request.Input.SpaceId != "" {
		spaceId = &request.Input.SpaceId
	}

	dbCtx, dbSpan := otel.Tracer("sqlc").Start(request.Context, "Insert permission set")
	permissionSet, err := sqlc.Queries.CreatePermissionSet(dbCtx, sqldb.CreatePermissionSetParams{
		ID:             *permissionSetId,
		OrganizationID: request.Input.OrgId,
		SpaceID:        spaceId,
		Name:           request.Input.Name,
		Permissions:    request.Input.Permissions,
		CreatedBy:      request.GetActorId(),
		UpdatedBy:      request.GetActorId(),
	})

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
		PermissionSet: &permissionSet,
	}

	return &output, nil
}
