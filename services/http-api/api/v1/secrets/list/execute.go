package secretsList

import (
	"fmt"
	"sort"

	"github.com/gofiber/fiber/v2"
	gruenthttp "github.com/ongruent/gruent/libs/http"
	"github.com/ongruent/gruent/services/http-api/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type Secret struct {
	Name string `json:"name"`
}

func execute(request *gruenthttp.Request[Input]) (*Output, error) {
	request.Span.AddEvent("execute")
	request.Span.SetAttributes(
		attribute.String("request.Input.SpaceId", request.Input.SpaceId),
	)

	// List existing space secrets
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

	spaceSecretNames := []Secret{}
	for key := range spaceSecrets {
		spaceSecretNames = append(spaceSecretNames, Secret{Name: key})
	}
	sort.Slice(spaceSecretNames, func(i, j int) bool {
		return spaceSecretNames[i].Name < spaceSecretNames[j].Name
	})

	output := &Output{
		Secrets: spaceSecretNames,
	}

	return output, nil
}
