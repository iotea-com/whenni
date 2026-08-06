package healthcheck

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ongruent/gruent/services/http-api/api/healthcheck/get"
)

func Register(app *fiber.App) {
	api := app.Group("/healthcheck")

	api.Get("/", get.Handler)
}
