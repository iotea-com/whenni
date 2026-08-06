package organizationsMembersAdd

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	gruenthttp "github.com/ongruent/gruent/libs/http"
	"github.com/ongruent/gruent/services/http-api/services/sqlc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func execute(request *gruenthttp.Request[Input]) (*Output, error) {
	request.Span.AddEvent("execute")
	request.Span.SetAttributes(
		attribute.String("request.Input.OrgId", request.Input.OrgId),
		attribute.String("request.Input.UserId", request.Input.UserId),
	)

	// Get organization's default permission set
	dbCtx, dbSpan := otel.Tracer("sqlc").Start(request.Context, "Get default permission set")
	permissionSets, err := sqlc.Queries.ListOrganizationPermissionSets(dbCtx, request.Input.OrgId)

	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error getting the organization's default permission set: %s", err)),
		)
		dbSpan.End()
		return nil, fiber.NewError(fiber.StatusInternalServerError, "could not get the organization's default permission set")
	}

	dbSpan.End()

	// Upsert member into organization
	defaultPermissionSetId := ""
	for _, permissionSet := range permissionSets {
		if permissionSet.Name == "Default" && permissionSet.SpaceID == nil {
			defaultPermissionSetId = permissionSet.ID
			break
		}
	}
	if defaultPermissionSetId == "" {
		request.Span.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", "default permission set not found"),
		)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "could not get the organization's default permission set")
	}

	dbCtx, dbSpan = otel.Tracer("sqlc").Start(request.Context, "Upsert member into organization")
	_, err = sqlc.Queries.CreateOrganizationMember(dbCtx, request.Input.OrgId, request.Input.UserId, "MEMBER", defaultPermissionSetId)

	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error upserting member into organization with ID %s in the database: %s", request.Input.OrgId, err)),
		)
		dbSpan.End()
		return nil, fiber.NewError(fiber.StatusInternalServerError, "could not upsert member into organization")
	}

	dbSpan.End()

	return nil, nil
}
