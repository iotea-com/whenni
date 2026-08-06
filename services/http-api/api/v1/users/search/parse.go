package search

import (
	"github.com/gofiber/fiber/v2"
	gruenthttp "github.com/ongruent/gruent/libs/http"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func parse(ctx *fiber.Ctx) (*gruenthttp.Request[Input], error) {
	requestSpan := trace.SpanFromContext(ctx.UserContext())
	requestSpan.AddEvent("parse")

	query := ctx.Query("q")
	orgId := ctx.Query("orgId")
	requestSpan.SetAttributes(
		attribute.String("query.q", query),
		attribute.String("query.orgId", orgId),
	)

	request := &gruenthttp.Request[Input]{
		Span:         requestSpan,
		FiberContext: ctx,
		Context:      ctx.UserContext(),
		Input: Input{
			Query: query,
			OrgId: orgId,
		},
	}

	return request, nil
}
