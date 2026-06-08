package permissionsRemove

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
		attribute.String("request.Input.Org", request.Input.OrgId),
		attribute.String("request.Input.PermissionSetId", request.Input.PermissionSetId),
	)

	// Check if permission set is the default permission set
	dbCtx, dbSpan := otel.Tracer("sqlc").Start(request.Context, "Check if permission set is the default permission set")
	permissionSet, err := sqlc.Queries.GetPermissionSet(dbCtx, request.Input.PermissionSetId)

	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error deleting permission set from the database: %s", err)),
		)
		dbSpan.End()
		return nil, fiber.NewError(fiber.StatusInternalServerError)
	}

	dbSpan.End()

	if permissionSet.Name == "Default" {
		errMessage := "Cannot delete the 'Default' permission set. This permission set is required for users that are not assigned a specific permission set."
		request.Span.SetAttributes(
			attribute.String("error.type", "input_validation"),
			attribute.String("error.message", errMessage),
		)
		errResponse := ioteahttp.NewErrorResponse([]any{errMessage})
		responseJson, _ := errResponse.MarshalJson()
		return nil, fiber.NewError(fiber.StatusBadRequest, string(responseJson))
	}

	// Delete permission set from the database
	dbCtx, dbSpan = otel.Tracer("sqlc").Start(request.Context, "Delete permission set")
	err = sqlc.Queries.DeletePermissionSet(dbCtx, request.Input.PermissionSetId)

	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error deleting permission set from the database: %s", err)),
		)
		dbSpan.End()
		return nil, fiber.NewError(fiber.StatusInternalServerError)
	}

	dbSpan.End()

	return nil, nil
}
