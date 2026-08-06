package channelsStatus

import (
	"github.com/gofiber/fiber/v2"
	gruenthttp "github.com/ongruent/gruent/libs/http"
)

type ResponseBody struct {
	Status string `json:"status"`
}

func respond(request *gruenthttp.Request[Input], output *Output) {
	request.Span.AddEvent("respond")

	responseBody := ResponseBody{
		Status: output.Status,
	}

	response := gruenthttp.NewGetResponse(responseBody)
	request.FiberContext.Status(fiber.StatusOK).JSON(response)
}
