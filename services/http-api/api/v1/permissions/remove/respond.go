package permissionsRemove

import (
	"github.com/gofiber/fiber/v2"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
)

func respond(request *ioteahttp.Request[Input], _ *Output) {
	request.Span.AddEvent("respond")

	response := ioteahttp.NewDeleteResponse()
	request.FiberContext.Status(fiber.StatusOK).JSON(response)
}
