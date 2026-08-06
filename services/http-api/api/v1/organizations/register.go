package organizations

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ongruent/gruent/services/http-api/api/middleware"
	organizationsCreate "github.com/ongruent/gruent/services/http-api/api/v1/organizations/create"
	organizationsDelete "github.com/ongruent/gruent/services/http-api/api/v1/organizations/delete"
	organizationsGet "github.com/ongruent/gruent/services/http-api/api/v1/organizations/get"
	organizationsMembers "github.com/ongruent/gruent/services/http-api/api/v1/organizations/members"
	organizationsUpdate "github.com/ongruent/gruent/services/http-api/api/v1/organizations/update"
)

func Register(app fiber.Router) {
	organizations := app.Group("/organizations")

	organizations.Use(middleware.ValidateBearerToken)
	organizations.Use(middleware.UserRateLimiter)

	// Members
	organizationsMembers.Register(organizations)

	// Organizations
	organizations.Post("/", organizationsCreate.Handler)
	organizations.Get("/:orgId", organizationsGet.Handler)
	organizations.Put("/:orgId", organizationsUpdate.Handler)
	organizations.Delete("/:orgId", organizationsDelete.Handler)
}
