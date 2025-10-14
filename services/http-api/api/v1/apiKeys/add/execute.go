package apiKeysAdd

import (
	"github.com/gofiber/fiber/v2"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	"github.com/iotea-com/iotea/libs/id"
	"github.com/iotea-com/iotea/prisma/db"
	"github.com/iotea-com/iotea/services/http-api/services/prisma"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func execute(request *ioteahttp.Request[Input]) (*Output, error) {
	request.Span.AddEvent("execute")
	request.Span.SetAttributes(
		attribute.String("request.Input.OrgId", request.Input.OrgId),
		attribute.String("request.Input.SpaceId", request.Input.SpaceId),
		attribute.String("request.Input.Name", request.Input.Name),
		attribute.String("request.Input.PermissionSetId", request.Input.PermissionSetId),
	)

	// Create API key ID
	apiKeyId, err := id.Generator.NewApiKeyId()
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "api_key_generation"),
			attribute.String("error.message", err.Error()),
		)
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Determine if we're creating an organization or space API key
	apiKeyType := "organization"
	if request.Input.SpaceId != "" {
		apiKeyType = "space"
	}

	// Store API key in the database
	var name *string
	if request.Input.Name != "" {
		name = &request.Input.Name
	}

	output := Output{
		ApiKey: nil,
	}

	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Insert API key")
	if apiKeyType == "space" {
		apiKey, err := prisma.Client.APIKey.CreateOne(
			db.APIKey.ID.Set(*apiKeyId),
			db.APIKey.CreatedBy.Set(request.GetActorId()),
			db.APIKey.Organization.Link(
				db.Organization.ID.Equals(request.Input.OrgId),
			),
			db.APIKey.Space.Link(
				db.Space.ID.Equals(request.Input.SpaceId),
			),
			db.APIKey.SpacePermissionSet.Link(
				db.PermissionSet.ID.Equals(request.Input.PermissionSetId),
			),
			db.APIKey.Name.SetIfPresent(name),
		).Exec(dbCtx)

		if err != nil {
			dbSpan.SetAttributes(
				attribute.String("error.type", "database"),
				attribute.String("error.message", err.Error()),
			)
			dbSpan.End()
			return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		output.ApiKey = apiKey
	} else {
		apiKey, err := prisma.Client.APIKey.CreateOne(
			db.APIKey.ID.Set(*apiKeyId),
			db.APIKey.CreatedBy.Set(request.GetActorId()),
			db.APIKey.Organization.Link(
				db.Organization.ID.Equals(request.Input.OrgId),
			),
			db.APIKey.OrganizationPermissionSet.Link(
				db.PermissionSet.ID.Equals(request.Input.PermissionSetId),
			),
			db.APIKey.Name.SetIfPresent(name),
		).Exec(dbCtx)

		if err != nil {
			dbSpan.SetAttributes(
				attribute.String("error.type", "database"),
				attribute.String("error.message", err.Error()),
			)
			dbSpan.End()
			return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		output.ApiKey = apiKey
	}

	dbSpan.SetAttributes(
		attribute.String("database.apiKey.id", output.ApiKey.ID),
	)
	dbSpan.End()

	return &output, nil
}
