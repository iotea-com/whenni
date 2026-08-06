package secretsCreate

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	gruenthttp "github.com/ongruent/gruent/libs/http"
	"github.com/ongruent/gruent/services/http-api/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func execute(request *gruenthttp.Request[Input]) (*Output, error) {
	request.Span.AddEvent("execute")
	request.Span.SetAttributes(
		attribute.String("request.Input.Name", request.Input.Name),
		// attribute.String("request.Input.Value", request.Input.Value), // sensitive information - must not be logged
		attribute.String("request.Input.SpaceId", request.Input.SpaceId),
	)

	// Get existing space secrets
	secretCtx, secretSpan := otel.Tracer("secretsClient").Start(request.Context, "List space secrets")
	if config.SecretsClient == nil {
		secretSpan.SetAttributes(
			attribute.String("error.type", "secrets"),
			attribute.String("error.message", "secrets client not initialized"),
		)
		secretSpan.End()

		errorResponse := gruenthttp.NewErrorResponse([]any{"secrets client not initialized"})
		responseJson, _ := errorResponse.MarshalJson()
		return nil, fiber.NewError(fiber.StatusBadRequest, string(responseJson))
	}

	secretsPath := fmt.Sprintf("spaces/%s", request.Input.SpaceId)
	spaceSecrets, err := config.SecretsClient.ReadKeyValue(secretCtx, secretsPath)
	if err != nil {
		if err.Error() != "secret not found" {
			secretSpan.SetAttributes(
				attribute.String("error.type", "secrets"),
				attribute.String("error.message", fmt.Sprintf("error listing space secrets: %s", err)),
			)
			secretSpan.End()

			errorResponse := gruenthttp.NewErrorResponse([]any{err.Error()})
			responseJson, _ := errorResponse.MarshalJson()
			return nil, fiber.NewError(fiber.StatusBadRequest, string(responseJson))
		}
	}

	secretSpan.End()

	// Check if secret already exists
	if _, ok := spaceSecrets[request.Input.Name]; ok {
		errorResponse := gruenthttp.NewErrorResponse([]any{"secret already exists"})
		responseJson, _ := errorResponse.MarshalJson()
		return nil, fiber.NewError(fiber.StatusBadRequest, string(responseJson))
	}

	// Append secret to space secrets
	secretCtx, secretSpan = otel.Tracer("secretsClient").Start(request.Context, "Create secret")
	newSecrets := map[string]any{}
	if spaceSecrets != nil {
		newSecrets = spaceSecrets
	}
	newSecrets[request.Input.Name] = request.Input.Value

	// Write the new space secret
	err = config.SecretsClient.WriteKeyValue(secretCtx, secretsPath, newSecrets)
	if err != nil {
		secretSpan.SetAttributes(
			attribute.String("error.type", "secrets"),
			attribute.String("error.message", fmt.Sprintf("error creating secret: %s", err)),
		)
		secretSpan.End()

		errorResponse := gruenthttp.NewErrorResponse([]any{err.Error()})
		responseJson, _ := errorResponse.MarshalJson()
		return nil, fiber.NewError(fiber.StatusBadRequest, string(responseJson))
	}

	secretSpan.End()

	return nil, nil
}
