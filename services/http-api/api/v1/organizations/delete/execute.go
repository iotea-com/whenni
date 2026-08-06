package organizationsDelete

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
	)

	orgId := request.Input.OrgId

	// Get spaces count for this organization from the database
	dbCtx, dbSpan := otel.Tracer("sqlc").Start(request.Context, "Count spaces")
	spacesCount, err := sqlc.Queries.CountSpacesByOrganization(dbCtx, orgId)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error counting spaces from the database: %s", err)),
		)
		dbSpan.End()
		return nil, err
	}
	dbSpan.End()

	// Do not delete space if there are spaces linked to the organization
	if spacesCount > 0 {
		request.Span.SetAttributes(
			attribute.String("error.type", "context_validation"),
			attribute.String("error.message", "delete organization failed since there are spaces linked to the organization"),
		)
		response := gruenthttp.NewErrorResponse([]any{"delete all spaces in this organization before deleting the organization"})
		responseJson, _ := response.MarshalJson()
		return nil, fiber.NewError(fiber.StatusConflict, string(responseJson))
	}

	// Delete api keys
	dbCtx, dbSpan = otel.Tracer("sqlc").Start(request.Context, "Delete API keys")
	apiKeys, err := sqlc.Queries.ListApiKeysByOrganization(dbCtx, orgId)
	if err == nil {
		for _, apiKey := range apiKeys {
			if err = sqlc.Queries.DeleteApiKey(dbCtx, apiKey.ID); err != nil {
				break
			}
		}
	}
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
	dbCtx, dbSpan = otel.Tracer("sqlc").Start(request.Context, "Delete member records")
	members, err := sqlc.Queries.ListOrganizationMembers(dbCtx, orgId)
	if err == nil {
		for _, member := range members {
			if err = sqlc.Queries.DeleteOrganizationMember(dbCtx, member.ID); err != nil {
				break
			}
		}
	}
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
	dbCtx, dbSpan = otel.Tracer("sqlc").Start(request.Context, "Delete permission sets")
	permissionSets, err := sqlc.Queries.ListOrganizationPermissionSets(dbCtx, orgId)
	if err == nil {
		for _, permissionSet := range permissionSets {
			if err = sqlc.Queries.DeletePermissionSet(dbCtx, permissionSet.ID); err != nil {
				break
			}
		}
	}
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
	dbCtx, dbSpan = otel.Tracer("sqlc").Start(request.Context, "Delete organization")
	err = sqlc.Queries.DeleteOrganization(dbCtx, orgId)
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
