package permissionsUpdate

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	gruenthttp "github.com/ongruent/gruent/libs/http"
	gruentpermissions "github.com/ongruent/gruent/libs/http/permissions"
	"github.com/ongruent/gruent/services/http-api/config"
	"github.com/ongruent/gruent/services/http-api/services/sqlc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func contextValidate(request *gruenthttp.Request[Input]) error {
	request.Span.AddEvent("contextValidate")

	authorizeRequestParams := gruenthttp.AuthorizeRequestParams{
		BearerToken: request.Input.BearerToken,
		JwtSecret:   config.VaultConf.JwtSecret,
		ScopeId:     request.Input.OrgId,
		Namespace:   gruentpermissions.NamespacePermissions,
		Action:      gruentpermissions.ActionUpdate,
	}

	if request.Input.SpaceId != "" {
		authorizeRequestParams.ScopeId = request.Input.SpaceId
	}

	err := request.Authorize(authorizeRequestParams)
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "authorize"),
			attribute.String("error.message", err.Error()),
			attribute.String("context_validation.status", "fail"),
		)
		return err
	}

	// Get the current permission set
	dbCtx, dbSpan := otel.Tracer("sqlc").Start(request.Context, "Get permission set")
	permissionSet, dbErr := sqlc.Queries.GetPermissionSet(dbCtx, request.Input.PermissionSetId)

	if dbErr != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error updating permission set in the database: %s", dbErr)),
		)
		dbSpan.End()
		return fiber.NewError(fiber.StatusInternalServerError, dbErr.Error())
	}

	// Check if the permission set is the Default permission set
	if permissionSet.Name == "Default" && request.Input.Name != "Default" {
		errMessage := "Cannot change the name of the 'Default' permission set."
		request.Span.SetAttributes(
			attribute.String("error.type", "context_validation"),
			attribute.String("error.message", errMessage),
		)
		errResponse := gruenthttp.NewErrorResponse([]any{errMessage})
		responseJson, _ := errResponse.MarshalJson()
		return fiber.NewError(fiber.StatusBadRequest, string(responseJson))
	}

	if permissionSet.Name != "Default" && request.Input.Name == "Default" {
		errMessage := "Cannot change the name of the permission set to 'Default'."
		request.Span.SetAttributes(
			attribute.String("error.type", "context_validation"),
			attribute.String("error.message", errMessage),
		)
		errResponse := gruenthttp.NewErrorResponse([]any{errMessage})
		responseJson, _ := errResponse.MarshalJson()
		return fiber.NewError(fiber.StatusBadRequest, string(responseJson))
	}

	request.Span.SetAttributes(attribute.String("context_validation.status", "pass"))
	return nil
}
