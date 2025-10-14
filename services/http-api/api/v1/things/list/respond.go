package thingsList

import (
	"github.com/gofiber/fiber/v2"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
)

func respond(request *ioteahttp.Request[Input], output *Output) {
	request.Span.AddEvent("respond")

	response := ioteahttp.NewListResponse(output.Things, &output.Page, &output.TotalPages, &output.TotalResults, &output.ResultsPerPage)
	request.FiberContext.Status(fiber.StatusOK).JSON(response)
}
