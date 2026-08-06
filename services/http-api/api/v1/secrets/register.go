package secrets

import (
	"github.com/gofiber/fiber/v2"

	"github.com/ongruent/gruent/services/http-api/api/middleware"
	secretsCreate "github.com/ongruent/gruent/services/http-api/api/v1/secrets/create"
	secretsDelete "github.com/ongruent/gruent/services/http-api/api/v1/secrets/delete"
	secretsList "github.com/ongruent/gruent/services/http-api/api/v1/secrets/list"
	secretsUpdate "github.com/ongruent/gruent/services/http-api/api/v1/secrets/update"
)

func Register(app fiber.Router) {
	secrets := app.Group("/secrets")

	secrets.Use(middleware.ValidateBearerToken)
	secrets.Use(middleware.UserRateLimiter)

	secretsByName := secrets.Group("/:secretName")

	secrets.Get("/", secretsList.Handler)
	secrets.Post("/", secretsCreate.Handler)
	secretsByName.Delete("/", secretsDelete.Handler)
	secretsByName.Put("/", secretsUpdate.Handler)
}
