package refresh

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	gruenthttp "github.com/ongruent/gruent/libs/http"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type RequestBody struct {
	RefreshToken string `json:"refreshToken"`
}

func parse(ctx *fiber.Ctx) (*gruenthttp.Request[Input], error) {
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

	requestSpan.SetAttributes(
		attribute.String("requestBody.refreshToken", requestBody.RefreshToken),
	)

	request := &gruenthttp.Request[Input]{
		Span:         requestSpan,
		FiberContext: ctx,
		Context:      ctx.UserContext(),
		Input: Input{
			RefreshToken: requestBody.RefreshToken,
		},
	}

	return request, nil
}
