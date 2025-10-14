package spacesUpdate

import (
	"encoding/json"
	"net/http"

	"github.com/gofiber/fiber/v2"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	ioteahttputil "github.com/iotea-com/iotea/libs/http/util"
	"github.com/iotea-com/iotea/prisma/db"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func parse(ctx *fiber.Ctx) (*ioteahttp.Request[Input], error) {
	type RequestBody struct {
		Space *db.SpaceModel `json:"space"`
	}

	requestSpan := trace.SpanFromContext(ctx.UserContext())
	requestSpan.AddEvent("parse")

	orgId := ctx.Query("orgId")
	spaceId := ctx.Params("spaceId")
	requestSpan.SetAttributes(
		attribute.String("params.orgId", orgId),
		attribute.String("params.spaceId", spaceId),
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
		return nil, fiber.NewError(http.StatusUnprocessableEntity, err.Error())
	}

	request := &ioteahttp.Request[Input]{
		Span:         requestSpan,
		FiberContext: ctx,
		Context:      ctx.UserContext(),
		Input: Input{
			BearerToken: *bearerToken,
			Space:       *requestBody.Space,
			OrgId:       orgId,
		},
	}

	return request, nil
}
