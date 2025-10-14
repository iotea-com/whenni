package tagsList

import (
	"github.com/gofiber/fiber/v2"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
)

func respond(request *ioteahttp.Request[Input], output *Output) {
	request.Span.AddEvent("respond")

	response := ioteahttp.NewGetResponse(output.Tags)
	request.FiberContext.Status(fiber.StatusOK).JSON(response)
}
