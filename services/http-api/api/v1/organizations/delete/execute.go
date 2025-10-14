package organizationsDelete

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	"github.com/iotea-com/iotea/prisma/db"
	"github.com/iotea-com/iotea/services/http-api/services/prisma"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func execute(request *ioteahttp.Request[Input]) (*Output, error) {
	request.Span.AddEvent("execute")
	request.Span.SetAttributes(
		attribute.String("request.Input.OrgId", request.Input.OrgId),
	)

	orgId := request.Input.OrgId

	// Get spaces for this organization from the database
	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Get spaces")
	spaces, err := prisma.Client.Space.FindMany(
		db.Space.OrganizationID.Equals(orgId),
	).Exec(dbCtx)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error getting spaces from the database: %s", err)),
		)
		dbSpan.End()
		return nil, err
	}
	dbSpan.End()

	// Do not delete space if there are spaces linked to the organization
	if len(spaces) > 0 {
		request.Span.SetAttributes(
			attribute.String("error.type", "context_validation"),
			attribute.String("error.message", "delete organization failed since there are spaces linked to the organization"),
		)
		response := ioteahttp.NewErrorResponse([]any{"delete all spaces in this organization before deleting the organization"})
		responseJson, _ := response.MarshalJson()
		return nil, fiber.NewError(fiber.StatusConflict, string(responseJson))
	}

	// Delete api keys
	dbCtx, dbSpan = otel.Tracer("prisma").Start(request.Context, "Delete API keys")
	_, err = prisma.Client.APIKey.FindMany(
		db.APIKey.OrganizationID.Equals(orgId),
	).Delete().Exec(dbCtx)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error deleting API keys from the database: %s", err)),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	// Delete all member records
	dbCtx, dbSpan = otel.Tracer("prisma").Start(request.Context, "Delete member records")
	_, err = prisma.Client.OrganizationMember.FindMany(
		db.OrganizationMember.OrganizationID.Equals(orgId),
	).Delete().Exec(dbCtx)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error deleting member records from the database: %s", err)),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	// Delete permission sets
	dbCtx, dbSpan = otel.Tracer("prisma").Start(request.Context, "Delete permission sets")
	_, err = prisma.Client.PermissionSet.FindMany(
		db.PermissionSet.OrganizationID.Equals(orgId),
	).Delete().Exec(dbCtx)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error deleting permission sets from the database: %s", err)),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	// Delete organization
	dbCtx, dbSpan = otel.Tracer("prisma").Start(request.Context, "Delete organization")
	_, err = prisma.Client.Organization.FindUnique(
		db.Organization.ID.Equals(orgId),
	).Delete().Exec(dbCtx)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error deleting organization from the database: %s", err)),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	output := Output{}

	return &output, nil
}
