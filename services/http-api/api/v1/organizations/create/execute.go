package organizationsCreate

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	ioteapermissions "github.com/iotea-com/iotea/libs/http/permissions"
	"github.com/iotea-com/iotea/libs/id"
	"github.com/iotea-com/iotea/prisma/db"
	"github.com/iotea-com/iotea/services/http-api/services/prisma"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func execute(request *ioteahttp.Request[Input]) (*Output, error) {
	request.Span.AddEvent("execute")
	request.Span.SetAttributes(
		attribute.String("request.Input.Name", request.Input.Name),
		attribute.String("request.ActorId", request.GetActorId()),
	)

	// Create organization ID
	orgId, err := id.Generator.NewOrganizationId()
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "id_generator"),
			attribute.String("error.message", fmt.Sprintf("error generating a new Organization ID: %s", err)),
		)
		errResponse := ioteahttp.NewErrorResponse([]any{"Could not create the organization. Please try again or contact support if the issue persists."})
		marshalledErrResponse, _ := errResponse.MarshalJson()
		return nil, fiber.NewError(fiber.StatusConflict, string(marshalledErrResponse))
	}

	// Get actor email
	user, err := prisma.Client.User.FindUnique(
		db.User.ID.Equals(request.GetActorId()),
	).Exec(request.Context)
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error getting user from the database: %s", err)),
		)
		return nil, err
	}

	email, ok := user.Email()
	if !ok || email == "" {
		request.Span.SetAttributes(
			attribute.String("error.type", "validation"),
			attribute.String("error.message", "user has no email address"),
		)
		errResponse := ioteahttp.NewErrorResponse([]any{"User account must have an email address to create an organization."})
		marshalledErrResponse, _ := errResponse.MarshalJson()
		return nil, fiber.NewError(fiber.StatusBadRequest, string(marshalledErrResponse))
	}

	// Create organization
	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Insert organization")
	organization, err := prisma.Client.Organization.CreateOne(
		db.Organization.ID.Set(*orgId),
		db.Organization.Name.Set(request.Input.Name),
		db.Organization.CreatedBy.Set(request.GetActorId()),
		db.Organization.UpdatedBy.Set(request.GetActorId()),
	).Exec(dbCtx)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error inserting organization into the database: %s", err)),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	// Create default permission set
	defaultPermissionSetId, err := id.Generator.NewPermissionSetId()
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "id_generator"),
			attribute.String("error.message", fmt.Sprintf("error generating a new Default Permission Set ID: %s", err)),
		)
		return nil, err
	}

	dbCtx, dbSpan = otel.Tracer("prisma").Start(request.Context, "Insert default permission set")
	_, err = prisma.Client.PermissionSet.CreateOne(
		db.PermissionSet.ID.Set(*defaultPermissionSetId),
		db.PermissionSet.Name.Set("Default"),
		db.PermissionSet.CreatedBy.Set(request.GetActorId()),
		db.PermissionSet.UpdatedBy.Set(request.GetActorId()),
		db.PermissionSet.Organization.Link(
			db.Organization.ID.Equals(organization.ID),
		),
		db.PermissionSet.Permissions.Set(ioteapermissions.DefaultMemberOrganizationPermissions),
	).Exec(dbCtx)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error inserting default permission set into the database: %s", err)),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	// Create initial member
	dbCtx, dbSpan = otel.Tracer("prisma").Start(request.Context, "Insert initial organization member")
	_, err = prisma.Client.OrganizationMember.CreateOne(
		db.OrganizationMember.Organization.Link(
			db.Organization.ID.Equals(organization.ID),
		),
		db.OrganizationMember.User.Link(
			db.User.ID.Equals(request.GetActorId()),
		),
		db.OrganizationMember.OrganizationPermissionSet.Link(
			db.PermissionSet.ID.Equals(*defaultPermissionSetId),
		),
		db.OrganizationMember.Role.Set(db.OrganizationRoleAdmin),
	).Exec(dbCtx)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error inserting initial organization member into the database: %s", err)),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	output := Output{
		Organization: organization,
	}

	return &output, nil
}
