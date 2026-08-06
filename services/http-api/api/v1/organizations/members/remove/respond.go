package organizationsMembersRemove

import (
	"github.com/gofiber/fiber/v2"
	gruenthttp "github.com/ongruent/gruent/libs/http"
)

func respond(request *gruenthttp.Request[Input], _ *Output) {
	request.Span.AddEvent("respond")

	response := gruenthttp.NewDeleteResponse()
	request.FiberContext.Status(fiber.StatusOK).JSON(response)
}
