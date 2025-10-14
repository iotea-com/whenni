package trigger

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	apitest "github.com/iotea-com/iotea/services/http-api/api/test"
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

		apitest.AssertRouteExists(t, routes, "GET", "/trigger/http/:channelId")
		apitest.AssertRouteExists(t, routes, "POST", "/trigger/http/:channelId")
		apitest.AssertRouteExists(t, routes, "PUT", "/trigger/http/:channelId")
		apitest.AssertRouteExists(t, routes, "DELETE", "/trigger/http/:channelId")
		apitest.AssertRouteExists(t, routes, "PATCH", "/trigger/http/:channelId")
		apitest.AssertRouteExists(t, routes, "HEAD", "/trigger/http/:channelId")
	})
}
