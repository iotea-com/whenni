package environments

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

		apitest.AssertRouteExists(t, routes, "GET", "/environments/")
		apitest.AssertRouteExists(t, routes, "POST", "/environments/")
		apitest.AssertRouteExists(t, routes, "GET", "/environments/:environmentId/ssh")
		apitest.AssertRouteExists(t, routes, "GET", "/environments/:environmentId/status")
		apitest.AssertRouteExists(t, routes, "POST", "/environments/:environmentId/start")
		apitest.AssertRouteExists(t, routes, "POST", "/environments/:environmentId/stop")
		apitest.AssertRouteExists(t, routes, "DELETE", "/environments/:environmentId")
	})
}
