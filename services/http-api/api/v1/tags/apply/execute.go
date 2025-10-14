package tagsApply

import (
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
		attribute.String("request.Input.TagId", request.Input.TagId),
		attribute.String("request.Input.SubjectId", request.Input.SubjectId),
		attribute.String("request.Input.SpaceId", request.Input.SpaceId),
	)

	// Create initial db arguments
	var subjectLinkArg db.AppliedTagSetParam = nil

	// Determine the type of subject
	idType, err := id.Parse(request.Input.SubjectId)
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "id_parse"),
			attribute.String("error.message", fmt.Sprintf("error parsing subject ID: %s", err)),
		)
		return nil, err
	}

	switch idType {
	case id.IdTypeThing:
		subjectLinkArg = db.AppliedTag.Thing.Link(db.Thing.ID.Equals(request.Input.SubjectId))
	case id.IdTypeChannel:
		subjectLinkArg = db.AppliedTag.Channel.Link(db.Channel.ID.Equals(request.Input.SubjectId))
	case id.IdTypeModel:
		subjectLinkArg = db.AppliedTag.Model.Link(db.Model.ID.Equals(request.Input.SubjectId))
	default:
		request.Span.SetAttributes(
			attribute.String("error.type", "id_parse"),
			attribute.String("error.message", fmt.Sprintf("invalid subject ID type: %s", idType)),
		)
		return nil, fmt.Errorf("invalid subject ID type for tag application: %s", idType)
	}

	// Apply tag to the subject
	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Apply tag")
	_, err = prisma.Client.AppliedTag.CreateOne(
		db.AppliedTag.Space.Link(db.Space.ID.Equals(request.Input.SpaceId)),
		db.AppliedTag.Tag.Link(db.Tag.ID.Equals(request.Input.TagId)),
		db.AppliedTag.CreatedBy.Set(request.GetActorId()),
		db.AppliedTag.UpdatedBy.Set(request.GetActorId()),
		subjectLinkArg,
	).Exec(dbCtx)
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
