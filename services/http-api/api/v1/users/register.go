package users

import (
	"github.com/gofiber/fiber/v2"

	"github.com/ongruent/gruent/services/http-api/api/middleware"
	"github.com/ongruent/gruent/services/http-api/api/v1/users/get"
	"github.com/ongruent/gruent/services/http-api/api/v1/users/search"
)

func Register(app fiber.Router) {
	users := app.Group("users")

	users.Use(middleware.ValidateBearerToken)
	users.Use(middleware.UserRateLimiter)

	users.Get("/search", search.Handler)
	users.Get("/:userId", get.Handler)
}
