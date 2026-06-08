package tagsApply

import (
	"fmt"

	sqldb "github.com/iotea-com/iotea/db/sqlc"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	"github.com/iotea-com/iotea/libs/id"
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

	// Build the subject-specific link fields for app.applied_tags.
	params := sqldb.ApplyTagParams{
		SpaceID:   request.Input.SpaceId,
		TagID:     request.Input.TagId,
		CreatedBy: request.GetActorId(),
		UpdatedBy: request.GetActorId(),
	}

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
		params.ThingID = &request.Input.SubjectId
	case id.IdTypeChannel:
		params.ChannelID = &request.Input.SubjectId
	case id.IdTypeModel:
		params.ModelID = &request.Input.SubjectId
	default:
		request.Span.SetAttributes(
			attribute.String("error.type", "id_parse"),
			attribute.String("error.message", fmt.Sprintf("invalid subject ID type: %s", idType)),
		)
		return nil, fmt.Errorf("invalid subject ID type for tag application: %s", idType)
	}

	// Apply tag to the subject
	dbCtx, dbSpan := otel.Tracer("sqlc").Start(request.Context, "Apply tag")
	_, err = sqlc.Queries.ApplyTag(dbCtx, params)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error applying tag in the database: %s", err)),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	return nil, nil
}
