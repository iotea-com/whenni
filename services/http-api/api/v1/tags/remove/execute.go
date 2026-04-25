package tagsRemove

import (
	"fmt"

	ioteahttp "github.com/iotea-com/iotea/libs/http"
	"github.com/iotea-com/iotea/services/http-api/services/sqlc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func execute(request *ioteahttp.Request[Input]) (*Output, error) {
	request.Span.AddEvent("execute")
	request.Span.SetAttributes(
		attribute.String("request.Input.TagId", request.Input.TagId),
		attribute.String("request.Input.SubjectId", request.Input.SubjectId),
		attribute.String("request.Input.SpaceId", request.Input.SpaceId),
	)

	// Delete applied tag from the database
	dbCtx, dbSpan := otel.Tracer("sqlc").Start(request.Context, "Delete applied tag")
	err := sqlc.Queries.RemoveAppliedTagBySubject(dbCtx, request.Input.TagId, &request.Input.SubjectId)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error removing applied tag from the database: %s", err)),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	return nil, nil
}
