package auth

import (
	"github.com/gofiber/fiber/v2"

	// "github.com/iotea-com/iotea/libs/http/middleware"
	"github.com/iotea-com/iotea/services/http-api/api/middleware"
	magicLinkVerify "github.com/iotea-com/iotea/services/http-api/api/v1/auth/magicLink/verify"
	"github.com/iotea-com/iotea/services/http-api/api/v1/auth/mqtt"
	"github.com/iotea-com/iotea/services/http-api/api/v1/auth/refresh"
	"github.com/iotea-com/iotea/services/http-api/api/v1/auth/signin"
	"github.com/iotea-com/iotea/services/http-api/api/v1/auth/signup"
	"github.com/iotea-com/iotea/services/http-api/api/v1/auth/updatePassword"
)

func Register(app fiber.Router) {
	auth := app.Group("auth")

	// auth.Use(middleware.IpRateLimiter) // TODO: Add IP whitelist so that services in the same network can freely access the MQTT route

	auth.Post("/mqtt", mqtt.Handler)
	auth.Post("/signin", signin.Handler)
	auth.Post("/signup", signup.Handler)
	auth.Post("/refresh", refresh.Handler)
	auth.Post("/magic-link/verify", magicLinkVerify.Handler)

	password := auth.Group("password")
	password.Use(middleware.ValidateBearerToken)
	password.Use(middleware.UserRateLimiter)
	password.Patch("/update", updatePassword.Handler)
}
