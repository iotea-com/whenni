package policies

import (
	"github.com/gofiber/fiber/v2"

	"github.com/ongruent/gruent/services/http-api/api/middleware"
	policiesGet "github.com/ongruent/gruent/services/http-api/api/v1/policies/get"
	policiesList "github.com/ongruent/gruent/services/http-api/api/v1/policies/list"
	policiesUpdate "github.com/ongruent/gruent/services/http-api/api/v1/policies/update"
)

func Register(app fiber.Router) {
	policies := app.Group("policies")

	policies.Use(middleware.ValidateBearerToken)
	policies.Use(middleware.UserRateLimiter)

	policies.Get("/", policiesList.Handler)
	policies.Get("/:certificateId", policiesGet.Handler)
	policies.Put("/:certificateId", policiesUpdate.Handler)
}
