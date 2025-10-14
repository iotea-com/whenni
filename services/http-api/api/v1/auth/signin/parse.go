package signin

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type RequestBody struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	Method     string `json:"method"`
	AppUrl     string `json:"appUrl"`
	RedirectTo string `json:"redirect"`
}

func parse(ctx *fiber.Ctx) (*ioteahttp.Request[Input], error) {
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
		attribute.String("request.Input.Email", requestBody.Email),
		attribute.String("request.Input.Method", requestBody.Method),
	)

	request := &ioteahttp.Request[Input]{
		Span:         requestSpan,
		FiberContext: ctx,
		Context:      ctx.UserContext(),
		Input: Input{
			Email:      requestBody.Email,
			Password:   requestBody.Password,
			Method:     requestBody.Method,
			AppUrl:     requestBody.AppUrl,
			RedirectTo: requestBody.RedirectTo,
		},
	}

	return request, nil
}
