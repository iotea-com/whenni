package users

import (
	"github.com/gofiber/fiber/v2"

	"github.com/iotea-com/iotea/services/http-api/api/middleware"
	"github.com/iotea-com/iotea/services/http-api/api/v1/users/get"
	"github.com/iotea-com/iotea/services/http-api/api/v1/users/search"
)

func Register(app fiber.Router) {
	users := app.Group("users")

	users.Use(middleware.ValidateBearerToken)
	users.Use(middleware.UserRateLimiter)

	users.Get("/search", search.Handler)
	users.Get("/:userId", get.Handler)
}
