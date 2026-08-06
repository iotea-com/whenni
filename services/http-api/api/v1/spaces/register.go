package spaces

import (
	"github.com/gofiber/fiber/v2"

	"github.com/ongruent/gruent/services/http-api/api/middleware"
	spacesCreate "github.com/ongruent/gruent/services/http-api/api/v1/spaces/create"
	spacesDelete "github.com/ongruent/gruent/services/http-api/api/v1/spaces/delete"
	spacesGet "github.com/ongruent/gruent/services/http-api/api/v1/spaces/get"
	spacesUpdate "github.com/ongruent/gruent/services/http-api/api/v1/spaces/update"
)

func Register(app fiber.Router) {
	spaces := app.Group("/spaces")

	spaces.Use(middleware.ValidateBearerToken)
	spaces.Use(middleware.UserRateLimiter)

	spacesById := spaces.Group("/:spaceId")
	spacesById.Use(middleware.ValidateSpaceById)

	spaces.Post("/", spacesCreate.Handler)
	spacesById.Get("/", spacesGet.Handler)
	spacesById.Delete("/", spacesDelete.Handler)
	spacesById.Put("/", spacesUpdate.Handler)
}
