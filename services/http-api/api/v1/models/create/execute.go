package modelsCreate

import (
	"encoding/json"
	"fmt"

	ioteahttp "github.com/iotea-com/iotea/libs/http"
	"github.com/iotea-com/iotea/libs/id"
	"github.com/iotea-com/iotea/prisma/db"
	"github.com/iotea-com/iotea/services/http-api/services/prisma"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func execute(request *ioteahttp.Request[Input]) (*Output, error) {
	request.Span.AddEvent("execute")
	request.Span.SetAttributes(
		attribute.String("request.Input.Name", request.Input.Name),
		attribute.String("request.Input.SpaceId", request.Input.SpaceId),
	)

	modelId, err := id.Generator.NewModelId()
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "id_generator"),
			attribute.String("error.message", fmt.Sprintf("error generating a new Model ID: %s", err)),
		)
		return nil, err
	}

	attributesBytes, err := json.Marshal(request.Input.Attributes)
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "parse_json"),
			attribute.String("error.message", fmt.Sprintf("error marshalling attributes: %s", err)),
		)
		return nil, err
	}

	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Insert model")
	model, err := prisma.Client.Model.CreateOne(
		db.Model.ID.Set(*modelId),
		db.Model.Name.Set(request.Input.Name),
		db.Model.Attributes.Set(attributesBytes),
		db.Model.CreatedBy.Set(request.GetActorId()),
		db.Model.UpdatedBy.Set(request.GetActorId()),
		db.Model.Space.Link(
			db.Space.ID.Equals(request.Input.SpaceId),
		),
	).Exec(dbCtx)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error inserting model into the database: %s", err)),
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
