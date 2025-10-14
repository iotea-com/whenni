package mqtt

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func parse(ctx *fiber.Ctx) (*ioteahttp.Request[Input], error) {
	type RequestBody struct {
		CertString string `json:"certString" validate:"required"`
		Topic      string `json:"topic" validate:"required"`
		Action     string `json:"action" validate:"required"`
	}

	requestSpan := trace.SpanFromContext(ctx.UserContext())
	requestSpan.AddEvent("parse")

	var requestBody RequestBody
	err := json.Unmarshal(ctx.Body(), &requestBody)
	if err != nil {
		requestSpan.SetAttributes(
			attribute.String("error.type", "parse_request_body"),
			attribute.String("error.message", fmt.Sprintf("invalid request body (%s): %s", err.Error(), ctx.Body())),
		)
		// respond with 200 and "result": "deny" so the request is blocked by the broker
		ctx.Status(200).JSON(ResponseBody{
			Result: "deny",
			Errors: []string{
				fmt.Sprintf("invalid request body (%s): %s", err.Error(), ctx.Body()),
			},
		})
		return nil, nil
	}

	request := &ioteahttp.Request[Input]{
		Span:         requestSpan,
		FiberContext: ctx,
		Context:      ctx.UserContext(),
		Input:        Input(requestBody),
	}

	return request, nil
}
