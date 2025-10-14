package organizationsMembersChangeRole

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	ioteapermissions "github.com/iotea-com/iotea/libs/http/permissions"
	"github.com/iotea-com/iotea/prisma/db"
	"github.com/iotea-com/iotea/services/http-api/config"
	"github.com/iotea-com/iotea/services/http-api/services/prisma"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func contextValidate(request *ioteahttp.Request[Input]) error {
	request.Span.AddEvent("contextValidate")

	authorizeRequestParams := ioteahttp.AuthorizeRequestParams{
		BearerToken:  request.Input.BearerToken,
		PrismaClient: prisma.Client,
		JwtSecret:    config.VaultConf.JwtSecret,
		ScopeId:      request.Input.OrgId,
		Namespace:    ioteapermissions.NamespaceMembers,
		Action:       ioteapermissions.ActionUpdate,
		EnforceAdmin: true, // Actor must be a user and an admin to change a member's role
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

	// Get admins count from the database
	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Count admins")
	var adminsCountResponse []struct {
		Count db.RawString `json:"admins_count"`
	}
	dbErr := prisma.Client.Prisma.QueryRaw(`SELECT count(*) as admins_count FROM app.organization_members WHERE organization_members."organizationId" = $1 AND organization_members."role" = 'ADMIN'`, request.Input.OrgId).Exec(dbCtx, &adminsCountResponse)
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error counting admins in the database: %s", err)),
		)
		dbSpan.End()

		return dbErr
	}

	adminsCount, dbErr := strconv.ParseInt(string(adminsCountResponse[0].Count), 10, 16)
	if dbErr != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("could not parse count response to int - %#v: %s", adminsCountResponse, err)),
		)
		dbSpan.End()

		return dbErr
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
