package apiKeysList

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

	page := ctx.QueryInt("page", 1)
	resultsPerPage := ctx.QueryInt("resultsPerPage", 10)

	orgId := ctx.Query("orgId")
	spaceId := ctx.Query("spaceId")
	requestSpan.SetAttributes(
		attribute.String("query.spaceId", spaceId),
		attribute.String("query.orgId", orgId),
		attribute.Int64("query.page", int64(page)),
		attribute.Int64("query.resultsPerPage", int64(resultsPerPage)),
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
			BearerToken:    *bearerToken,
			OrgId:          orgId,
			SpaceId:        spaceId,
			Page:           page,
			ResultsPerPage: resultsPerPage,
		},
	}

	return request, nil
}
