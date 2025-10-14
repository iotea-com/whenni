package organizationsMembersChangeRole

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
		attribute.String("request.Input.UserId", request.Input.UserId),
	)

	// Set member's role in organization
	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Set member's role in organization")
	_, err := prisma.Client.OrganizationMember.FindUnique(
		db.OrganizationMember.OrganizationIDUserID(
			db.OrganizationMember.OrganizationID.Equals(request.Input.OrgId),
			db.OrganizationMember.UserID.Equals(request.Input.UserId),
		),
	).Update(
		db.OrganizationMember.Role.Set(db.OrganizationRole(request.Input.Role)),
	).Exec(dbCtx)

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
