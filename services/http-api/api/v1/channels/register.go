package channels

import (
	"github.com/gofiber/fiber/v2"

	"github.com/ongruent/gruent/services/http-api/api/middleware"
	channelsCreate "github.com/ongruent/gruent/services/http-api/api/v1/channels/create"
	channelsDelete "github.com/ongruent/gruent/services/http-api/api/v1/channels/delete"
	channelExecutions "github.com/ongruent/gruent/services/http-api/api/v1/channels/executions"
	channelsGet "github.com/ongruent/gruent/services/http-api/api/v1/channels/get"
	channelsList "github.com/ongruent/gruent/services/http-api/api/v1/channels/list"
	channelsPublish "github.com/ongruent/gruent/services/http-api/api/v1/channels/publish"
	channelsStatus "github.com/ongruent/gruent/services/http-api/api/v1/channels/status"
	channelsUnpublish "github.com/ongruent/gruent/services/http-api/api/v1/channels/unpublish"
	channelsUpdateConfig "github.com/ongruent/gruent/services/http-api/api/v1/channels/updateConfig"
	channelsValidate "github.com/ongruent/gruent/services/http-api/api/v1/channels/validate"
)

func Register(app fiber.Router) {
	channels := app.Group("channels")

	channels.Use(middleware.ValidateBearerToken)
	channels.Use(middleware.UserRateLimiter)

	// Executions
	channelExecutions.Register(channels)

	// Channels
	channels.Post("/", channelsCreate.Handler)
	channels.Get("/", channelsList.Handler)
	channels.Patch("/validate", channelsValidate.Handler)

	channelsById := channels.Group("/:channelId")
	channelsById.Delete("/", channelsDelete.Handler)
	channelsById.Get("/", channelsGet.Handler)
	channelsById.Patch("/publish", channelsPublish.Handler)
	channelsById.Patch("/unpublish", channelsUnpublish.Handler)
	channelsById.Patch("/config", channelsUpdateConfig.Handler)
	channelsById.Get("/status", channelsStatus.Handler)
}
