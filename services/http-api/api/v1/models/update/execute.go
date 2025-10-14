package modelsUpdate

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	"github.com/iotea-com/iotea/prisma/db"
	"github.com/iotea-com/iotea/services/http-api/services/prisma"
	"github.com/iotea-com/iotea/services/http-api/util"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func execute(request *ioteahttp.Request[Input]) (*Output, error) {
	request.Span.AddEvent("execute")
	request.Span.SetAttributes(
		attribute.String("request.Input.ModelId", request.Input.ModelId),
		attribute.String("request.Input.SpaceId", request.Input.SpaceId),
	)

	now := util.GetCurrentTime()

	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Update struct in space")

	attributesBytes, err := json.Marshal(request.Input.attributes)
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "parse_json"),
			attribute.String("error.message", fmt.Sprintf("error marshalling attributes: %s", err)),
		)
		return nil, fiber.NewError(fiber.StatusUnprocessableEntity, "error marshalling attributes")
	}

	model, err := prisma.Client.Model.FindUnique(
		db.Model.ID.Equals(request.Input.ModelId),
	).Update(
		db.Model.Name.Set(request.Input.Model.Name),
		db.Model.Attributes.Set(attributesBytes),
		db.Model.UpdatedBy.Set(request.GetActorId()),
		db.Model.UpdatedAt.Set(now),
	).Exec(dbCtx)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error updating model in the database: %s", err)),
		)
		dbSpan.End()

		return nil, err
	}
	dbSpan.End()

	output := Output{
		Model: model,
	}

	return &output, nil
}
