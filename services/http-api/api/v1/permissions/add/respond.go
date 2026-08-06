package permissionsAdd

import (
	"github.com/gofiber/fiber/v2"
	gruenthttp "github.com/ongruent/gruent/libs/http"
)

func respond(request *gruenthttp.Request[Input], output *Output) {
	request.Span.AddEvent("respond")

	response := gruenthttp.NewCreateResponse(output.PermissionSet)
	request.FiberContext.Status(fiber.StatusCreated).JSON(response)
}
