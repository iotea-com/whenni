package tagsRemove

import (
	"fmt"

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
		attribute.String("request.Input.SubjectId", request.Input.SubjectId),
		attribute.String("request.Input.SpaceId", request.Input.SpaceId),
	)

	// Delete applied tag from the database
	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Delete applied tag")
	_, err := prisma.Client.AppliedTag.FindMany(
		db.AppliedTag.TagID.Equals(request.Input.TagId),
		db.AppliedTag.Or(
			db.AppliedTag.ThingID.Equals(request.Input.SubjectId),
			db.AppliedTag.ChannelID.Equals(request.Input.SubjectId),
			db.AppliedTag.ModelID.Equals(request.Input.SubjectId),
		),
	).Delete().Exec(dbCtx)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error inserting space into the database: %s", err)),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	return nil, nil
}
