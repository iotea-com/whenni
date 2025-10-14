package environments

import (
	"github.com/gofiber/fiber/v2"

	"github.com/iotea-com/iotea/services/http-api/api/middleware"
	environmentsCreate "github.com/iotea-com/iotea/services/http-api/api/v1/environments/create"
	environmentsDelete "github.com/iotea-com/iotea/services/http-api/api/v1/environments/delete"
	environmentsList "github.com/iotea-com/iotea/services/http-api/api/v1/environments/list"
	environmentsSsh "github.com/iotea-com/iotea/services/http-api/api/v1/environments/ssh"
	environmentsStart "github.com/iotea-com/iotea/services/http-api/api/v1/environments/start"
	environmentsStatus "github.com/iotea-com/iotea/services/http-api/api/v1/environments/status"
	environmentsStop "github.com/iotea-com/iotea/services/http-api/api/v1/environments/stop"
)

func Register(app fiber.Router) {
	environments := app.Group("environments")

	environments.Use(middleware.ValidateBearerToken)
	environments.Use(middleware.UserRateLimiter)

	environments.Post("/", environmentsCreate.Handler)
	environments.Get("/", environmentsList.Handler)
	environments.Get("/:environmentId/ssh", environmentsSsh.Handler)
	environments.Get("/:environmentId/status", environmentsStatus.Handler)
	environments.Post("/:environmentId/start", environmentsStart.Handler)
	environments.Post("/:environmentId/stop", environmentsStop.Handler)
	environments.Delete("/:environmentId", environmentsDelete.Handler)
}
