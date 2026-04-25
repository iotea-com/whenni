package organizationsMembersChangeRole

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	ioteapermissions "github.com/iotea-com/iotea/libs/http/permissions"
	"github.com/iotea-com/iotea/services/http-api/config"
	"github.com/iotea-com/iotea/services/http-api/services/sqlc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func contextValidate(request *ioteahttp.Request[Input]) error {
	request.Span.AddEvent("contextValidate")

	authorizeRequestParams := ioteahttp.AuthorizeRequestParams{
		BearerToken:  request.Input.BearerToken,
		JwtSecret:    config.VaultConf.JwtSecret,
		ScopeId:      request.Input.OrgId,
		Namespace:    ioteapermissions.NamespaceMembers,
		Action:       ioteapermissions.ActionUpdate,
		EnforceAdmin: true, // Actor must be a user and an admin to change a member's role
	}

	authErr := request.Authorize(authorizeRequestParams)
	if authErr != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "authorize"),
			attribute.String("error.message", authErr.Error()),
			attribute.String("context_validation.status", "fail"),
		)
		return authErr
	}

	// Get admins count from the database
	dbCtx, dbSpan := otel.Tracer("sqlc").Start(request.Context, "Count admins")
	members, dbErr := sqlc.Queries.ListOrganizationMembers(dbCtx, request.Input.OrgId)
	if dbErr != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error counting admins in the database: %s", dbErr)),
		)
		dbSpan.End()

		return dbErr
	}

	adminsCount := 0
	for _, member := range members {
		if member.Role == "ADMIN" {
			adminsCount++
		}
	}

	// Check that there is at least one admin in the organization
	if adminsCount <= 1 && request.Input.Role == "MEMBER" {
		errMsg := "there must be at least one admin in the organization"
		request.Span.SetAttributes(
			attribute.String("error.type", "context_validation"),
			attribute.String("error.message", errMsg),
			attribute.String("context_validation.status", "fail"),
		)

		errorResponse := ioteahttp.NewErrorResponse([]any{errMsg})
		marshalledErrorResponse, _ := errorResponse.MarshalJson()

		return fiber.NewError(fiber.StatusBadRequest, string(marshalledErrorResponse))
	}

	request.Span.SetAttributes(attribute.String("context_validation.status", "pass"))
	return nil
}
