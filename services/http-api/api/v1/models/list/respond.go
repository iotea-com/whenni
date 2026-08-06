package modelsList

import (
	"github.com/gofiber/fiber/v2"
	gruenthttp "github.com/ongruent/gruent/libs/http"
)

func respond(request *gruenthttp.Request[Input], output *Output) {
	request.Span.AddEvent("respond")

	response := gruenthttp.NewListResponse(output.Models, &output.Page, &output.TotalPages, &output.TotalResults, &output.ResultsPerPage)
	request.FiberContext.Status(fiber.StatusOK).JSON(response)
}
