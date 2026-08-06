package permissions

import (
	"github.com/gofiber/fiber/v2"

	"github.com/ongruent/gruent/services/http-api/api/middleware"
	permissionsAdd "github.com/ongruent/gruent/services/http-api/api/v1/permissions/add"
	permissionsList "github.com/ongruent/gruent/services/http-api/api/v1/permissions/list"
	permissionsRemove "github.com/ongruent/gruent/services/http-api/api/v1/permissions/remove"
	permissionsUpdate "github.com/ongruent/gruent/services/http-api/api/v1/permissions/update"
)

func Register(app fiber.Router) {
	permissions := app.Group("permissions")

	permissions.Use(middleware.ValidateBearerToken)
	permissions.Use(middleware.UserRateLimiter)

	permissions.Get("/", permissionsList.Handler)
	permissions.Post("/", permissionsAdd.Handler)
	permissions.Put("/:permissionSetId", permissionsUpdate.Handler)
	permissions.Delete("/", permissionsRemove.Handler)
}
