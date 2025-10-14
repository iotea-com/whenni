package channelId

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/gofiber/fiber/v2"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	"go.opentelemetry.io/otel/attribute"
)

func respond(request *ioteahttp.Request[Input], output *Output) {
	request.Span.AddEvent("respond")

	responseBodyBytes, err := io.ReadAll(output.response.Body)
	if err != nil {
		errorMessage := fmt.Sprintf("could not read response body from HTTP server: %s", err.Error())
		request.Span.SetAttributes(
			attribute.String("error.type", "read_body"),
			attribute.String("error.message", errorMessage),
		)
		r := ioteahttp.IoteaApiResponse{
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

	var response ioteahttp.IoteaApiResponse
	switch output.response.Request.Method {
	case fiber.MethodPost:
		response = ioteahttp.NewCreateResponse(responseBody)
	case fiber.MethodPut:
		response = ioteahttp.NewUpdateResponse(responseBody, nil)
	case fiber.MethodDelete:
		response = ioteahttp.NewDeleteResponse()
	default:
		response = ioteahttp.NewGetResponse(responseBody)
	}

	request.FiberContext.Status(output.response.StatusCode).JSON(response)
}
