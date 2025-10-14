package models

import (
	"github.com/gofiber/fiber/v2"

	"github.com/iotea-com/iotea/services/http-api/api/middleware"
	modelsCreate "github.com/iotea-com/iotea/services/http-api/api/v1/models/create"
	modelsDelete "github.com/iotea-com/iotea/services/http-api/api/v1/models/delete"
	modelsGet "github.com/iotea-com/iotea/services/http-api/api/v1/models/get"
	modelsList "github.com/iotea-com/iotea/services/http-api/api/v1/models/list"
	modelsUpdate "github.com/iotea-com/iotea/services/http-api/api/v1/models/update"
)

func Register(app fiber.Router) {
	models := app.Group("models")

	models.Use(middleware.ValidateBearerToken)
	models.Use(middleware.UserRateLimiter)

	models.Get("/", modelsList.Handler)
	models.Post("/", modelsCreate.Handler)
	models.Get("/:modelId", modelsGet.Handler)
	models.Delete("/:modelId", modelsDelete.Handler)
	models.Put("/:modelId", modelsUpdate.Handler)
}
