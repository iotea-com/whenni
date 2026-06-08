package apiKeysRemove

import (
	"github.com/gofiber/fiber/v2"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	ioteapermissions "github.com/iotea-com/iotea/libs/http/permissions"
	"github.com/iotea-com/iotea/services/http-api/config"
	"github.com/iotea-com/iotea/services/http-api/services/sqlc"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func contextValidate(request *ioteahttp.Request[Input]) error {
	request.Span.AddEvent("contextValidate")

	authorizeRequestParams := ioteahttp.AuthorizeRequestParams{
		BearerToken: request.Input.BearerToken,
		JwtSecret:   config.VaultConf.JwtSecret,
		ScopeId:     request.Input.OrgId,
		Namespace:   ioteapermissions.NamespaceOrganizationApiKeys,
		Action:      ioteapermissions.ActionDelete,
	}

	if request.Input.SpaceId != "" {
		authorizeRequestParams.ScopeId = request.Input.SpaceId
		authorizeRequestParams.Namespace = ioteapermissions.NamespaceSpaceApiKeys
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

	// Get the API key from the database
	dbCtx, dbSpan := otel.Tracer("sqlc").Start(request.Context, "Delete API key")
	apiKey, apiKeyFindErr := sqlc.Queries.GetApiKey(dbCtx, request.Input.ApiKeyId)
	if apiKeyFindErr != nil {
		if apiKeyFindErr == pgx.ErrNoRows {
			dbSpan.SetAttributes(
				attribute.String("error.type", "database"),
				attribute.String("error.message", apiKeyFindErr.Error()),
				attribute.String("context_validation.status", "fail"),
			)
			dbSpan.End()
			return fiber.NewError(fiber.StatusNotFound)
		}

		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", apiKeyFindErr.Error()),
			attribute.String("context_validation.status", "fail"),
		)
		dbSpan.End()
		return fiber.NewError(fiber.StatusInternalServerError)
	}

	dbSpan.End()

	// Verify that the resource is a space API key if a space ID is provided
	if apiKey.SpaceID != nil && request.Input.SpaceId == "" {
		errMsg := "attempted to delete a space API key without providing a space ID"
		request.Span.SetAttributes(
			attribute.String("error.type", "validate_space_api_key"),
			attribute.String("error.message", errMsg),
			attribute.String("context_validation.status", "fail"),
		)
		errResponse := ioteahttp.NewErrorResponse([]any{errMsg})
		marshalledErrResponse, _ := errResponse.MarshalJson()
		return fiber.NewError(fiber.StatusConflict, string(marshalledErrResponse))
	}

	// Verify that the resource is an organization API key if no space ID is provided
	if apiKey.SpaceID == nil && request.Input.SpaceId != "" {
		errMsg := "attempted to delete an organization API key, but a space ID was provided"
		request.Span.SetAttributes(
			attribute.String("error.type", "validate_organization_api_key"),
			attribute.String("error.message", errMsg),
			attribute.String("context_validation.status", "fail"),
		)
		errResponse := ioteahttp.NewErrorResponse([]any{errMsg})
		marshalledErrResponse, _ := errResponse.MarshalJson()
		return fiber.NewError(fiber.StatusConflict, string(marshalledErrResponse))
	}

	request.Span.SetAttributes(attribute.String("context_validation.status", "pass"))
	return nil
}
