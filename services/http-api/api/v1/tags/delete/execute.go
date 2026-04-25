package tagsDelete

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	"github.com/iotea-com/iotea/services/http-api/services/sqlc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func execute(request *ioteahttp.Request[Input]) (*Output, error) {
	request.Span.AddEvent("execute")
	request.Span.SetAttributes(
		attribute.String("request.Input.TagId", request.Input.TagId),
		attribute.String("request.Input.SpaceId", request.Input.SpaceId),
	)

	// Check if there are any applied tags
	dbCtx, dbSpan := otel.Tracer("sqlc").Start(request.Context, "Check if there are any applied tags")
	appliedTagCount, err := sqlc.Queries.CountAppliedTagsByTag(dbCtx, request.Input.TagId)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error checking if tag is applied: %s", err)),
		)
		dbSpan.End()

		return nil, err
	}

	dbSpan.End()

	if appliedTagCount > 0 {
		errMessage := fmt.Sprintf("tag is used in %d thing(s), model(s), or channel(s)", appliedTagCount)
		request.Span.SetAttributes(
			attribute.String("error.type", "context_validation"),
			attribute.String("error.message", errMessage),
		)

		errorResponse := ioteahttp.NewErrorResponse([]any{
			errMessage,
		})
		errorResponseJson, _ := errorResponse.MarshalJson()
		return nil, fiber.NewError(fiber.StatusConflict, string(errorResponseJson))
	}

	return nil, nil
}
