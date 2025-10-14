package magicLinkVerify

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func parse(ctx *fiber.Ctx) (*ioteahttp.Request[Input], error) {
	type RequestBody struct {
		Token string `json:"token"`
	}

	requestSpan := trace.SpanFromContext(ctx.UserContext())
	requestSpan.AddEvent("parse")

	var requestBody RequestBody
	err := json.Unmarshal(ctx.Body(), &requestBody)
	if err != nil {
		requestSpan.SetAttributes(
			attribute.String("error.type", "parse_request_body"),
			attribute.String("error.message", err.Error()),
		)
		return nil, fiber.NewError(fiber.StatusUnprocessableEntity)
	}

	token := requestBody.Token
	requestSpan.SetAttributes(
		attribute.String("request.Input.VerificationToken", token),
	)

	request := &ioteahttp.Request[Input]{
		Span:         requestSpan,
		FiberContext: ctx,
		Context:      ctx.UserContext(),
		Input: Input{
			VerificationToken: token,
		},
	}

	return request, nil
}
