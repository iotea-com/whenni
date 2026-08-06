package secretsUpdate

import (
	"encoding/json"
	"net/http"

	"github.com/gofiber/fiber/v2"
	gruenthttp "github.com/ongruent/gruent/libs/http"
	gruenthttputil "github.com/ongruent/gruent/libs/http/util"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func parse(ctx *fiber.Ctx) (*gruenthttp.Request[Input], error) {
	type RequestBody struct {
		Value string `json:"value"`
	}

	requestSpan := trace.SpanFromContext(ctx.UserContext())
	requestSpan.AddEvent("parse")

	spaceId := ctx.Query("spaceId")
	secretName := ctx.Params("secretName")
	requestSpan.SetAttributes(
		attribute.String("query.spaceId", spaceId),
		attribute.String("params.secretName", secretName),
	)

	bearerToken, err := gruenthttputil.ParseBearerToken(ctx.Get("Authorization"))
	if err != nil {
		requestSpan.SetAttributes(
			attribute.String("error.type", "parse_bearer_token"),
			attribute.String("error.message", err.Error()),
		)
		return nil, err
	}

	var requestBody RequestBody
	err = json.Unmarshal(ctx.Body(), &requestBody)
	if err != nil {
		requestSpan.SetAttributes(
			attribute.String("error.type", "parse_request_body"),
			attribute.String("error.message", err.Error()),
		)
		return nil, fiber.NewError(http.StatusBadRequest)
	}

	request := &gruenthttp.Request[Input]{
		Span:         requestSpan,
		FiberContext: ctx,
		Context:      ctx.UserContext(),
		Input: Input{
			BearerToken: *bearerToken,
			Name:        secretName,
			Value:       requestBody.Value,
			SpaceId:     spaceId,
		},
	}

	return request, nil
}
