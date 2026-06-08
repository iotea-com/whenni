package trigger

import (
	"github.com/gofiber/fiber/v2"
	"github.com/iotea-com/iotea/libs/legacy/engine/environment"
	"github.com/iotea-com/iotea/services/http-api/api/middleware"
	"github.com/iotea-com/iotea/services/http-api/api/v1/trigger/http/channelId"
	"github.com/iotea-com/iotea/services/http-api/config"
)

func Register(app fiber.Router) {
	trigger := app.Group("/trigger")

	// We do performance testing on local and development environments,
	// so we want to let as many requests through as possible.
	if config.Env == environment.Production {
		trigger.Use(middleware.IpRateLimiter)
	}

	// HTTP trigger
	trigger.Get("/http/:channelId", channelId.Handler)
	trigger.Post("/http/:channelId", channelId.Handler)
	trigger.Put("/http/:channelId", channelId.Handler)
	trigger.Delete("/http/:channelId", channelId.Handler)
	trigger.Patch("/http/:channelId", channelId.Handler)
	trigger.Head("/http/:channelId", channelId.Handler)
}
