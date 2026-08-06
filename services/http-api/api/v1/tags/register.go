package tags

import (
	"github.com/gofiber/fiber/v2"

	"github.com/ongruent/gruent/services/http-api/api/middleware"
	tagsApply "github.com/ongruent/gruent/services/http-api/api/v1/tags/apply"
	tagsCreate "github.com/ongruent/gruent/services/http-api/api/v1/tags/create"
	tagsDelete "github.com/ongruent/gruent/services/http-api/api/v1/tags/delete"
	tagsList "github.com/ongruent/gruent/services/http-api/api/v1/tags/list"
	tagsRemove "github.com/ongruent/gruent/services/http-api/api/v1/tags/remove"
)

func Register(app fiber.Router) {
	tags := app.Group("/tags")

	tags.Use(middleware.ValidateBearerToken)
	tags.Use(middleware.UserRateLimiter)

	tagsByName := tags.Group("/:tagId")

	tags.Get("/", tagsList.Handler)
	tags.Post("/", tagsCreate.Handler)
	tagsByName.Delete("/", tagsDelete.Handler)
	tagsByName.Patch("/apply", tagsApply.Handler)
	tagsByName.Patch("/remove", tagsRemove.Handler)
}
