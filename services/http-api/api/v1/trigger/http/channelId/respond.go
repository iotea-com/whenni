package channelId

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/gofiber/fiber/v2"
	gruenthttp "github.com/ongruent/gruent/libs/http"
	"go.opentelemetry.io/otel/attribute"
)

func respond(request *gruenthttp.Request[Input], output *Output) {
	request.Span.AddEvent("respond")

	responseBodyBytes, err := io.ReadAll(output.response.Body)
	if err != nil {
		errorMessage := fmt.Sprintf("could not read response body from HTTP server: %s", err.Error())
		request.Span.SetAttributes(
			attribute.String("error.type", "read_body"),
			attribute.String("error.message", errorMessage),
		)
		r := gruenthttp.GruentApiResponse{
			Data:   nil,
			Errors: []any{errorMessage},
		}
		request.FiberContext.Status(fiber.StatusInternalServerError).JSON(r)
		return
	}

	var responseBody interface{} = nil
	if len(responseBodyBytes) > 0 {
		err = json.Unmarshal(responseBodyBytes, &responseBody)
		if err != nil {
			request.Span.SetAttributes(
				attribute.String("error.type", "unmarshal_response_body"),
				attribute.String("error.message", err.Error()),
			)
		}

		request.Span.SetAttributes(
			attribute.Int("response.status", output.response.StatusCode),
		)
	}

	var response gruenthttp.GruentApiResponse
	switch output.response.Request.Method {
	case fiber.MethodPost:
		response = gruenthttp.NewCreateResponse(responseBody)
	case fiber.MethodPut:
		response = gruenthttp.NewUpdateResponse(responseBody, nil)
	case fiber.MethodDelete:
		response = gruenthttp.NewDeleteResponse()
	default:
		response = gruenthttp.NewGetResponse(responseBody)
	}

	request.FiberContext.Status(output.response.StatusCode).JSON(response)
}
