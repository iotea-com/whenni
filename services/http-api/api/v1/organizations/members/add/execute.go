package organizationsMembersAdd

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

	// Get organization's default permission set
	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Get default permission set")
	defaultPermissionSet, err := prisma.Client.PermissionSet.FindFirst(
		db.PermissionSet.And(
			db.PermissionSet.OrganizationID.Equals(request.Input.OrgId),
			db.PermissionSet.SpaceID.IsNull(),
			db.PermissionSet.Name.Equals("Default"),
		),
	).Exec(dbCtx)

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
	dbCtx, dbSpan = otel.Tracer("prisma").Start(request.Context, "Upsert member into organization")
	_, err = prisma.Client.OrganizationMember.CreateOne(
		db.OrganizationMember.Organization.Link(
			db.Organization.ID.Equals(request.Input.OrgId),
		),
		db.OrganizationMember.User.Link(
			db.User.ID.Equals(request.Input.UserId),
		),
		db.OrganizationMember.OrganizationPermissionSet.Link(
			db.PermissionSet.ID.Equals(defaultPermissionSet.ID),
		),
		db.OrganizationMember.Role.Set(db.OrganizationRoleMember),
	).Exec(dbCtx)

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
