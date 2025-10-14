package things

import (
	"github.com/gofiber/fiber/v2"

	"github.com/iotea-com/iotea/services/http-api/api/middleware"
	thingsCreate "github.com/iotea-com/iotea/services/http-api/api/v1/things/create"
	thingsDelete "github.com/iotea-com/iotea/services/http-api/api/v1/things/delete"
	thingsGet "github.com/iotea-com/iotea/services/http-api/api/v1/things/get"
	thingsHealthcheck "github.com/iotea-com/iotea/services/http-api/api/v1/things/healthcheck"
	thingsList "github.com/iotea-com/iotea/services/http-api/api/v1/things/list"
	thingsUpdate "github.com/iotea-com/iotea/services/http-api/api/v1/things/update"
)

func Register(app fiber.Router) {
	things := app.Group("things")

	things.Use(middleware.ValidateBearerToken)
	things.Use(middleware.UserRateLimiter)

	things.Get("/", thingsList.Handler)
	things.Post("/", thingsCreate.Handler)
	things.Patch("/healthcheck", thingsHealthcheck.Handler)
	things.Get("/:thingId", thingsGet.Handler)
	things.Put("/:thingId", thingsUpdate.Handler)
	things.Delete("/:thingId", thingsDelete.Handler)
}
