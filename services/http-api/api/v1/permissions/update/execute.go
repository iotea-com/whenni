package permissionsUpdate

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	"github.com/iotea-com/iotea/services/http-api/services/sqlc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func execute(request *ioteahttp.Request[Input]) (*Output, error) {
	request.Span.AddEvent("execute")
	request.Span.SetAttributes(
		attribute.String("request.Input.Name", request.Input.Name),
		attribute.String("request.Input.OrgId", request.Input.OrgId),
		attribute.String("request.Input.PermissionSetId", request.Input.PermissionSetId),
		attribute.StringSlice("request.Input.ChannelId", request.Input.Permissions),
		attribute.String("request.Input.SpaceId", request.Input.SpaceId),
	)

	// Update the permission set
	var spaceId *string
	if request.Input.SpaceId != "" {
		spaceId = &request.Input.SpaceId
	}

	dbCtx, dbSpan := otel.Tracer("sqlc").Start(request.Context, "Update permission set")
	permissionSet, err := sqlc.Queries.UpdatePermissionSet(
		dbCtx,
		request.Input.PermissionSetId,
		request.Input.Name,
		request.Input.Permissions,
		request.GetActorId(),
		spaceId,
	)

	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error updating permission set in the database: %s", err)),
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
