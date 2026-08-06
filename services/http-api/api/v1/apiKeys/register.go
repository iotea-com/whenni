package apiKeys

import (
	"github.com/gofiber/fiber/v2"

	"github.com/ongruent/gruent/services/http-api/api/middleware"
	apiKeysAdd "github.com/ongruent/gruent/services/http-api/api/v1/apiKeys/add"
	apiKeysList "github.com/ongruent/gruent/services/http-api/api/v1/apiKeys/list"
	apiKeysRemove "github.com/ongruent/gruent/services/http-api/api/v1/apiKeys/remove"
)

func Register(app fiber.Router) {
	apiKeys := app.Group("api-keys")

	apiKeys.Use(middleware.ValidateBearerToken)
	apiKeys.Use(middleware.UserRateLimiter)

	apiKeys.Get("/", apiKeysList.Handler)
	apiKeys.Post("/", apiKeysAdd.Handler)
	apiKeys.Delete("/", apiKeysRemove.Handler)
}
