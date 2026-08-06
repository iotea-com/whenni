package channels

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

		apitest.AssertRouteExists(t, routes, "GET", "/channels/")
		apitest.AssertRouteExists(t, routes, "POST", "/channels/")
		apitest.AssertRouteExists(t, routes, "PATCH", "/channels/validate")

		apitest.AssertRouteExists(t, routes, "GET", "/channels/:channelId/")
		apitest.AssertRouteExists(t, routes, "DELETE", "/channels/:channelId/")
		apitest.AssertRouteExists(t, routes, "PATCH", "/channels/:channelId/publish")
		apitest.AssertRouteExists(t, routes, "PATCH", "/channels/:channelId/unpublish")
		apitest.AssertRouteExists(t, routes, "PATCH", "/channels/:channelId/config")
	})
}
