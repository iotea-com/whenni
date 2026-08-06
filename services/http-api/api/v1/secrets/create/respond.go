package secretsCreate

import (
	"github.com/gofiber/fiber/v2"
	gruenthttp "github.com/ongruent/gruent/libs/http"
)

func respond(request *gruenthttp.Request[Input], _ *Output) {
	request.Span.AddEvent("respond")

	response := gruenthttp.NewCreateResponse(nil)
	request.FiberContext.Status(fiber.StatusCreated).JSON(response)
}
