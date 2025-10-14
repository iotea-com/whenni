package tagsRemove

import (
	"github.com/gofiber/fiber/v2"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	ioteahttputil "github.com/iotea-com/iotea/libs/http/util"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func parse(ctx *fiber.Ctx) (*ioteahttp.Request[Input], error) {
	requestSpan := trace.SpanFromContext(ctx.UserContext())
	requestSpan.AddEvent("parse")

	tagId := ctx.Params("tagId")
	spaceId := ctx.Query("spaceId")
	subjectId := ctx.Query("subjectId")
	requestSpan.SetAttributes(
		attribute.String("params.tagId", tagId),
		attribute.String("query.spaceId", spaceId),
		attribute.String("query.subjectId", subjectId),
	)

	bearerToken, err := ioteahttputil.ParseBearerToken(ctx.Get("Authorization"))
	if err != nil {
		requestSpan.SetAttributes(
			attribute.String("error.type", "parse_bearer_token"),
			attribute.String("error.message", err.Error()),
		)
		return nil, err
	}

	request := &ioteahttp.Request[Input]{
		Span:         requestSpan,
		FiberContext: ctx,
		Context:      ctx.UserContext(),
		Input: Input{
			BearerToken: *bearerToken,
			TagId:       tagId,
			SubjectId:   subjectId,
			SpaceId:     spaceId,
		},
	}

	return request, nil
}
