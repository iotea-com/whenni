package channelsList

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	gruenthttp "github.com/ongruent/gruent/libs/http"
	gruenthttputil "github.com/ongruent/gruent/libs/http/util"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func parse(ctx *fiber.Ctx) (*gruenthttp.Request[Input], error) {
	requestSpan := trace.SpanFromContext(ctx.UserContext())
	requestSpan.AddEvent("parse")

	spaceId := ctx.Query("spaceId")
	filter := ctx.Query("q")
	tagFilterString := ctx.Query("tagFilter")
	page := ctx.QueryInt("page", 1)
	resultsPerPage := ctx.QueryInt("resultsPerPage", 10)
	requestSpan.SetAttributes(
		attribute.String("query.spaceId", spaceId),
		attribute.Int64("query.page", int64(page)),
		attribute.Int64("query.resultsPerPage", int64(resultsPerPage)),
		attribute.String("query.filter", filter),
		attribute.String("query.tagFilter", tagFilterString),
	)

	bearerToken, err := gruenthttputil.ParseBearerToken(ctx.Get("Authorization"))
	if err != nil {
		requestSpan.SetAttributes(
			attribute.String("error.type", "parse_bearer_token"),
			attribute.String("error.message", err.Error()),
		)
		return nil, err
	}

	tagFilter := []string{}
	if tagFilterString != "" {
		tagFilter = strings.Split(tagFilterString, ",")
	}

	request := &gruenthttp.Request[Input]{
		Span:         requestSpan,
		FiberContext: ctx,
		Context:      ctx.UserContext(),
		Input: Input{
			BearerToken:    *bearerToken,
			SpaceId:        spaceId,
			Page:           page,
			ResultsPerPage: resultsPerPage,
			Filter:         filter,
			TagFilter:      tagFilter,
		},
	}

	return request, nil
}
