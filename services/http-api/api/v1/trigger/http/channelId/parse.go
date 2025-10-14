package channelId

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func parse(ctx *fiber.Ctx) (*ioteahttp.Request[Input], error) {
	requestSpan := trace.SpanFromContext(ctx.UserContext())
	requestSpan.AddEvent("parse")

	// parse params
	channelId := ctx.Params("channelId")

	// Determine the Content-Type of the request
	contentType := ctx.Get("Content-Type")

	// parse body
	var requestBody map[string]any = nil
	var binaryData []byte

	// Handle JSON data
	if contentType == "application/json" {
		if len(ctx.Body()) > 0 {
			err := json.Unmarshal(ctx.Body(), &requestBody)
			if err != nil {
				requestSpan.SetAttributes(
					attribute.String("error.type", "parse_request_body"),
					attribute.String("error.message", err.Error()),
				)
				return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
			}
		}
	} else if contentType == "application/octet-stream" || contentType == "image/jpeg" || contentType == "image/png" || contentType == "application/pdf" {
		binaryData = ctx.Body() // Read binary data directly
	} else {
		return nil, fiber.NewError(fiber.StatusUnsupportedMediaType, "Unsupported Content-Type")
	}

	// instantiate request
	request := &ioteahttp.Request[Input]{
		Span:         requestSpan,
		FiberContext: ctx,
		Context:      ctx.UserContext(),
		Input: Input{
			ChannelId:         channelId,
			Method:            ctx.Method(),
			Body:              requestBody,
			BinaryData:        binaryData,
			BinaryContentType: contentType,
		},
	}

	return request, nil
}
