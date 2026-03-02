package modelsCreate

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	ioteahttputil "github.com/iotea-com/iotea/libs/http/util"
	"github.com/iotea-com/iotea/libs/legacy/engine/dependencies/models"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func parse(ctx *fiber.Ctx) (*ioteahttp.Request[Input], error) {
	type RequestBody struct {
		Name       string                      `json:"name"`
		Attributes map[string]models.Attribute `json:"attributes"`
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
		validationError := ioteahttp.NewErrorResponse([]any{err.Error()})
		responseJson, err := validationError.MarshalJson()
		if err != nil {
			return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		return nil, fiber.NewError(fiber.StatusBadRequest, string(responseJson))
	}

	request := &ioteahttp.Request[Input]{
		Span:         requestSpan,
		FiberContext: ctx,
		Context:      ctx.UserContext(),
		Input: Input{
			BearerToken: *bearerToken,
			SpaceId:     spaceId,
			Name:        requestBody.Name,
			Attributes:  requestBody.Attributes,
		},
	}

	return request, nil
}
