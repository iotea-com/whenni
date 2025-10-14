package tagsDelete

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	"github.com/iotea-com/iotea/prisma/db"
	"github.com/iotea-com/iotea/services/http-api/services/prisma"
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
	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Check if there are any applied tags")
	appliedTags, err := prisma.Client.AppliedTag.FindMany(
		db.AppliedTag.TagID.Equals(request.Input.TagId),
	).Exec(dbCtx)
	if err != nil && err.Error() != "ErrNotFound" {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error checking if secret is used in a thing: %s", err)),
		)
		dbSpan.End()

		return nil, err
	}

	dbSpan.End()

	if len(appliedTags) > 0 {
		errMessage := fmt.Sprintf("tag is used in %d thing(s), model(s), or channel(s)", len(appliedTags))
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
