package tagsCreate

import (
	"fmt"

	ioteahttp "github.com/iotea-com/iotea/libs/http"
	"github.com/iotea-com/iotea/libs/id"
	"github.com/iotea-com/iotea/services/http-api/services/sqlc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func execute(request *ioteahttp.Request[Input]) (*Output, error) {
	request.Span.AddEvent("execute")
	request.Span.SetAttributes(
		attribute.String("request.Input.Name", request.Input.Name),
		attribute.String("request.Input.SpaceId", request.Input.SpaceId),
	)

	// Create tag entry
	dbCtx, dbSpan := otel.Tracer("sqlc").Start(request.Context, "Insert tag")
	tagId, err := id.Generator.NewTagId()
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "id_generation"),
			attribute.String("error.message", fmt.Sprintf("error generating a new Tag ID: %s", err)),
		)
		return nil, err
	}

	tag, err := sqlc.Queries.CreateTag(
		dbCtx,
		*tagId,
		request.Input.Name,
		request.Input.SpaceId,
		request.GetActorId(),
		request.GetActorId(),
	)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error inserting tag into the database: %s", err)),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	output := &Output{
		Tag: tag,
	}

	return output, nil
}
