package spacesGet

import (
	"github.com/gofiber/fiber/v2"
	gruenthttp "github.com/ongruent/gruent/libs/http"
	gruenthttputil "github.com/ongruent/gruent/libs/http/util"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func parse(ctx *fiber.Ctx) (*gruenthttp.Request[Input], error) {
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

	request := &gruenthttp.Request[Input]{
		Span:         requestSpan,
		FiberContext: ctx,
		Context:      ctx.UserContext(),
		Input: Input{
			BearerToken: *bearerToken,
			OrgId:       orgId,
			SpaceId:     spaceId,
		},
	}

	return request, nil
}
