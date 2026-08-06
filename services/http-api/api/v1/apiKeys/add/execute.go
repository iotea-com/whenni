package apiKeysAdd

import (
	"github.com/gofiber/fiber/v2"
	sqldb "github.com/ongruent/gruent/db/sqlc"
	gruenthttp "github.com/ongruent/gruent/libs/http"
	"github.com/ongruent/gruent/libs/id"
	"github.com/ongruent/gruent/services/http-api/services/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func execute(request *gruenthttp.Request[Input]) (*Output, error) {
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

	dbCtx, dbSpan := otel.Tracer("sqlc").Start(request.Context, "Insert API key")
	if apiKeyType == "space" {
		apiKey, err := sqlc.Queries.CreateApiKey(dbCtx, sqldb.CreateApiKeyParams{
			ID:                          *apiKeyId,
			OrganizationID:              request.Input.OrgId,
			OrganizationPermissionSetID: nil,
			SpaceID:                     &request.Input.SpaceId,
			SpacePermissionSetID:        &request.Input.PermissionSetId,
			Name:                        name,
			CreatedBy:                   request.GetActorId(),
			ExpiresAt:                   pgtype.Timestamptz{},
		})

		if err != nil {
			dbSpan.SetAttributes(
				attribute.String("error.type", "database"),
				attribute.String("error.message", err.Error()),
			)
			dbSpan.End()
			return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		output.ApiKey = &apiKey
	} else {
		apiKey, err := sqlc.Queries.CreateApiKey(dbCtx, sqldb.CreateApiKeyParams{
			ID:                          *apiKeyId,
			OrganizationID:              request.Input.OrgId,
			OrganizationPermissionSetID: &request.Input.PermissionSetId,
			SpaceID:                     nil,
			SpacePermissionSetID:        nil,
			Name:                        name,
			CreatedBy:                   request.GetActorId(),
			ExpiresAt:                   pgtype.Timestamptz{},
		})

		if err != nil {
			dbSpan.SetAttributes(
				attribute.String("error.type", "database"),
				attribute.String("error.message", err.Error()),
			)
			dbSpan.End()
			return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		output.ApiKey = &apiKey
	}

	dbSpan.SetAttributes(
		attribute.String("database.apiKey.id", output.ApiKey.ID),
	)
	dbSpan.End()

	return &output, nil
}
