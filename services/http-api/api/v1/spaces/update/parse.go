package spacesUpdate

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
		Space *SpacePayload `json:"space"`
	}

	requestSpan := trace.SpanFromContext(ctx.UserContext())
	requestSpan.AddEvent("parse")

	orgId := ctx.Query("orgId")
	spaceId := ctx.Params("spaceId")
	requestSpan.SetAttributes(
		attribute.String("params.orgId", orgId),
		attribute.String("params.spaceId", spaceId),
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
		return nil, fiber.NewError(http.StatusUnprocessableEntity, err.Error())
	}
	if requestBody.Space == nil {
		return nil, fiber.NewError(http.StatusUnprocessableEntity, "space is required")
	}

	request := &gruenthttp.Request[Input]{
		Span:         requestSpan,
		FiberContext: ctx,
		Context:      ctx.UserContext(),
		Input: Input{
			BearerToken: *bearerToken,
			Space: SpacePayload{
				ID:   spaceId,
				Name: requestBody.Space.Name,
			},
			OrgId: orgId,
		},
	}

	return request, nil
}
