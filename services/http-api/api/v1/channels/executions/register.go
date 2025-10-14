package channelExecutions

import (
	"github.com/gofiber/fiber/v2"

	"github.com/iotea-com/iotea/services/http-api/api/v1/channels/executions/get"
	"github.com/iotea-com/iotea/services/http-api/api/v1/channels/executions/list"
)

func Register(app fiber.Router) {
	channelExecutions := app.Group("executions")

	channelExecutions.Get("/", list.Handler)
	channelExecutions.Get("/:executionId", get.Handler)
}
