package channelsStatus

import (
	"github.com/gofiber/fiber/v2"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
)

type ResponseBody struct {
	Status string `json:"status"`
}

func respond(request *ioteahttp.Request[Input], output *Output) {
	request.Span.AddEvent("respond")

	responseBody := ResponseBody{
		Status: output.Status.String(),
	}

	response := ioteahttp.NewGetResponse(responseBody)
	request.FiberContext.Status(fiber.StatusOK).JSON(response)
}
