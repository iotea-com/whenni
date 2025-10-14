package thingsCreate

import (
	"encoding/json"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/iotea-com/iotea/libs/engine/dependencies/things"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	ioteahttputil "github.com/iotea-com/iotea/libs/http/util"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func parse(ctx *fiber.Ctx) (*ioteahttp.Request[Input], error) {
	type RequestBody struct {
		Name       string               `json:"name"`
		Category   things.ThingCategory `json:"category"`
		Attributes map[string]any       `json:"attributes"`
	}

	requestSpan := trace.SpanFromContext(ctx.UserContext())
	requestSpan.AddEvent("parse")

	spaceId := ctx.Query("spaceId")
	requestSpan.SetAttributes(
		attribute.String("query.spaceId", spaceId),
	)

	bearerToken, err := ioteahttputil.ParseBearerToken(ctx.Get("Authorization"))
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
		return nil, fiber.NewError(http.StatusUnprocessableEntity)
	}

	request := &ioteahttp.Request[Input]{
		Span:         requestSpan,
		FiberContext: ctx,
		Context:      ctx.UserContext(),
		Input: Input{
			BearerToken: *bearerToken,
			SpaceId:     spaceId,
			Name:        requestBody.Name,
			Category:    requestBody.Category,
			Attributes:  requestBody.Attributes,
		},
	}

	return request, nil
}
