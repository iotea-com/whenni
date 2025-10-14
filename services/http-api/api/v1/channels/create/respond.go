package channelsCreate

import (
	"github.com/gofiber/fiber/v2"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
)

func respond(request *ioteahttp.Request[Input], output *Output) {
	request.Span.AddEvent("respond")

	response := ioteahttp.NewCreateResponse(output.Channel)
	request.FiberContext.Status(fiber.StatusCreated).JSON(response)
}
