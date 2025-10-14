package secretsUpdate

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	"github.com/iotea-com/iotea/services/http-api/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func execute(request *ioteahttp.Request[Input]) (*Output, error) {
	request.Span.AddEvent("execute")
	request.Span.SetAttributes(
		attribute.String("request.Input.Name", request.Input.Name),
		// attribute.String("request.Input.Value", request.Input.Value), // sensitive information - must not be logged
		attribute.String("request.Input.SpaceId", request.Input.SpaceId),
	)

	// Update an existing space secret
	secretCtx, secretSpan := otel.Tracer("secretsClient").Start(request.Context, "List space secrets")
	if config.SecretsClient == nil {
		secretSpan.SetAttributes(
			attribute.String("error.type", "secrets"),
			attribute.String("error.message", "secrets client not initialized"),
		)
		secretSpan.End()

		errorResponse := ioteahttp.NewErrorResponse([]any{"secrets client not initialized"})
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

			errorResponse := ioteahttp.NewErrorResponse([]any{err.Error()})
			responseJson, _ := errorResponse.MarshalJson()
			return nil, fiber.NewError(fiber.StatusBadRequest, string(responseJson))
		}
	}

	secretSpan.End()

	// Check if secret already exists
	if _, ok := spaceSecrets[request.Input.Name]; !ok {
		errorResponse := ioteahttp.NewErrorResponse([]any{"secret not found"})
		responseJson, _ := errorResponse.MarshalJson()
		return nil, fiber.NewError(fiber.StatusBadRequest, string(responseJson))
	}

	// Update the secret
	secretCtx, secretSpan = otel.Tracer("secretsClient").Start(request.Context, "Update secret")
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

		errorResponse := ioteahttp.NewErrorResponse([]any{err.Error()})
		responseJson, _ := errorResponse.MarshalJson()
		return nil, fiber.NewError(fiber.StatusBadRequest, string(responseJson))
	}

	secretSpan.End()

	return nil, nil
}
