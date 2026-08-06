package models

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	apitest "github.com/ongruent/gruent/services/http-api/api/test"
)

func TestRegister(t *testing.T) {
	t.Run("successfully registers", func(t *testing.T) {
		// given
		// create a new fiber app
		app := fiber.New()

		// when
		// register is called with a Fiber app
		Register(app)

		// then
		// assert it has the appropriate routes
		routes := app.GetRoutes()

		apitest.AssertRouteExists(t, routes, "GET", "/models/")
		apitest.AssertRouteExists(t, routes, "POST", "/models/")
		apitest.AssertRouteExists(t, routes, "GET", "/models/:modelId")
		apitest.AssertRouteExists(t, routes, "DELETE", "/models/:modelId")
		apitest.AssertRouteExists(t, routes, "PUT", "/models/:modelId")
	})
}
