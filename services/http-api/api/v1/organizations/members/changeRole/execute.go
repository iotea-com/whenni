package organizationsMembersChangeRole

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	sqldb "github.com/ongruent/gruent/db/sqlc"
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

	// Set member's role in organization
	dbCtx, dbSpan := otel.Tracer("sqlc").Start(request.Context, "Set member's role in organization")
	member, err := sqlc.Queries.GetOrganizationMember(dbCtx, request.Input.OrgId, request.Input.UserId)
	if err == nil {
		_, err = sqlc.Queries.UpdateOrganizationMemberRole(dbCtx, member.ID, sqldb.AppOrganizationRole(request.Input.Role))
	}

	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error setting member's role in organization with ID %s in the database: %s", request.Input.OrgId, err)),
		)
		dbSpan.End()
		return nil, fiber.NewError(fiber.StatusInternalServerError, "could not set member's role in organization")
	}

	dbSpan.End()

	return nil, nil
}
