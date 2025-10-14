package signup

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type RequestBody struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	AppUrl      string `json:"appUrl"`
	InviteToken string `json:"inviteToken"`
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
	)

	input := Input{
		Email:    requestBody.Email,
		Password: requestBody.Password,
		AppUrl:   requestBody.AppUrl,
	}

	if requestBody.InviteToken != "" {
		input.InviteToken = &requestBody.InviteToken
	}

	request := &ioteahttp.Request[Input]{
		Span:         requestSpan,
		FiberContext: ctx,
		Context:      ctx.UserContext(),
		Input:        input,
	}

	return request, nil
}
