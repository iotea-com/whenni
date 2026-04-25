package secretsDelete

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	sqldb "github.com/iotea-com/iotea/db/sqlc"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	"github.com/iotea-com/iotea/services/http-api/config"
	"github.com/iotea-com/iotea/services/http-api/services/sqlc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func execute(request *ioteahttp.Request[Input]) (*Output, error) {
	request.Span.AddEvent("execute")
	request.Span.SetAttributes(
		attribute.String("request.Input.Name", request.Input.Name),
		attribute.String("request.Input.SpaceId", request.Input.SpaceId),
	)

	// Check if secret is used in any things
	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Check if secret is used in things")
	things, err := sqlc.Queries.ListThings(dbCtx, request.Input.SpaceId)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error checking if secret is used in a thing: %s", err)),
		)
		dbSpan.End()

		return nil, err
	}

	dbSpan.End()

	thingsWithSecret := []sqldb.AppThing{}
	for _, thing := range things {
		var thingAttributes map[string]any
		err = json.Unmarshal(thing.Attributes, &thingAttributes)
		if err != nil {
			request.Span.SetAttributes(
				attribute.String("error.type", "parse_json"),
				attribute.String("error.message", fmt.Sprintf("could not read thing config: %s", err.Error())),
			)

			errorResponse := ioteahttp.NewErrorResponse([]any{
				fmt.Sprintf("could not read thing config: %s", err.Error()),
			})
			errorResponseJson, _ := errorResponse.MarshalJson()
			return nil, fiber.NewError(fiber.StatusBadRequest, string(errorResponseJson))
		}

		for _, value := range thingAttributes {
			if value == request.Input.Name {
				thingsWithSecret = append(thingsWithSecret, thing)
			}
		}
	}

	if len(thingsWithSecret) > 0 {
		errMessage := fmt.Sprintf("secret is used in %d thing(s)", len(thingsWithSecret))
		request.Span.SetAttributes(
			attribute.String("error.type", "context_validation"),
			attribute.String("error.message", errMessage),
		)

		errorResponse := ioteahttp.NewErrorResponse([]any{
			errMessage,
		})
		errorResponseJson, _ := errorResponse.MarshalJson()
		return nil, fiber.NewError(fiber.StatusConflict, string(errorResponseJson))
	}

	// Get existing space secrets
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

	// Remove secret from space secrets
	secretCtx, secretSpan = otel.Tracer("secretsClient").Start(request.Context, "Delete secret")
	newSecrets := map[string]any{}
	if spaceSecrets != nil {
		newSecrets = spaceSecrets
	}
	delete(newSecrets, request.Input.Name)

	// Write the new space secret
	if len(newSecrets) < 1 {
		err = config.SecretsClient.DeleteKeyValue(secretCtx, secretsPath)
		if err != nil {
			secretSpan.SetAttributes(
				attribute.String("error.type", "secrets"),
				attribute.String("error.message", fmt.Sprintf("error deleting secret: %s", err)),
			)
			secretSpan.End()

			errorResponse := ioteahttp.NewErrorResponse([]any{err.Error()})
			responseJson, _ := errorResponse.MarshalJson()
			return nil, fiber.NewError(fiber.StatusBadRequest, string(responseJson))
		}

		secretSpan.End()
	} else {
		err = config.SecretsClient.WriteKeyValue(secretCtx, secretsPath, newSecrets)
		if err != nil {
			secretSpan.SetAttributes(
				attribute.String("error.type", "secrets"),
				attribute.String("error.message", fmt.Sprintf("error deleting secret: %s", err)),
			)
			secretSpan.End()

			errorResponse := ioteahttp.NewErrorResponse([]any{err.Error()})
			responseJson, _ := errorResponse.MarshalJson()
			return nil, fiber.NewError(fiber.StatusBadRequest, string(responseJson))
		}

		secretSpan.End()
	}

	return nil, nil
}
