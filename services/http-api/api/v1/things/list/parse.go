package thingsList

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	ioteahttputil "github.com/iotea-com/iotea/libs/http/util"
	"github.com/iotea-com/iotea/libs/legacy/engine/dependencies/things"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func parse(ctx *fiber.Ctx) (*ioteahttp.Request[Input], error) {
	requestSpan := trace.SpanFromContext(ctx.UserContext())
	requestSpan.AddEvent("parse")

	page := ctx.QueryInt("page", 1)
	resultsPerPage := ctx.QueryInt("resultsPerPage", 10)

	searchFilter := ctx.Query("q")
	tagFilterString := ctx.Query("tagFilter")
	spaceId := ctx.Query("spaceId")
	bearerToken, err := ioteahttputil.ParseBearerToken(ctx.Get("Authorization"))
	if err != nil {
		requestSpan.SetAttributes(
			attribute.String("error.type", "parse_bearer_token"),
			attribute.String("error.message", err.Error()),
		)
		return nil, err
	}

	var thingCategory *things.ThingCategory
	category := ctx.Query("category")

	if category != "" {
		tc, ok := things.ParseCategoryString(category)
		if ok {
			thingCategory = tc
		} else {
			errResponse := ioteahttp.NewErrorResponse([]any{
				"Invalid category.",
			})
			responseJson, err := errResponse.MarshalJson()
			if err != nil {
				return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
			}
			return nil, fiber.NewError(fiber.StatusBadRequest, string(responseJson))
		}
	}

	tagFilter := []string{}
	if tagFilterString != "" {
		tagFilter = strings.Split(tagFilterString, ",")
	}

	request := &ioteahttp.Request[Input]{
		Span:         requestSpan,
		FiberContext: ctx,
		Context:      ctx.UserContext(),
		Input: Input{
			BearerToken:    *bearerToken,
			SpaceId:        spaceId,
			ThingCategory:  thingCategory,w
			Page:           page,
			ResultsPerPage: resultsPerPage,
			Filter:         searchFilter,
			TagFilter:      tagFilter,
		},
	}

	return request, nil
}
