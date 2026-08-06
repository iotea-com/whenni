package certificates

import (
	"github.com/gofiber/fiber/v2"

	"github.com/ongruent/gruent/services/http-api/api/middleware"
	certificatesGet "github.com/ongruent/gruent/services/http-api/api/v1/certificates/get"
)

func Register(app fiber.Router) {
	certificates := app.Group("certificates")

	certificates.Use(middleware.ValidateBearerToken)
	certificates.Use(middleware.UserRateLimiter)

	certificates.Get("/:certificateId", certificatesGet.Handler)
}
