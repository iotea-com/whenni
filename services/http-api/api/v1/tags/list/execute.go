package tagsList

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
		attribute.String("request.Input.SpaceId", request.Input.SpaceId),
		attribute.String("request.Input.Category", request.Input.Category),
	)

	// List tags in space (based on the filter)
	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "List tags in space")

	tags, err := prisma.Client.Tag.FindMany(
		db.Tag.SpaceID.Equals(request.Input.SpaceId),
	).With(
		db.Tag.Applications.Fetch().With(
			db.AppliedTag.Thing.Fetch(),
			db.AppliedTag.Channel.Fetch(),
			db.AppliedTag.Model.Fetch(),
		),
	).OrderBy(db.Tag.UpdatedAt.Order(db.DESC)).
		Exec(dbCtx)
	if err != nil {
		if err.Error() == "ErrNotFound" {
			request.Span.SetAttributes(
				attribute.String("error.type", "http_request"),
				attribute.String("error.message", fmt.Sprintf("no tags found in space with ID %s", request.Input.SpaceId)),
			)
			dbSpan.End()

			output := Output{
				Tags: nil,
			}

			return &output, nil
		}

		dbSpan.End()

		return nil, err
	}

	dbSpan.End()

	dbSpan.SetAttributes(
		attribute.String("response.Tags", fmt.Sprintf("%v", tags)),
	)

	// Group applied tags by tag
	tagMap := make(map[string]ExpandedTag)
	for _, tag := range tags {
		// Create new tag in map if necessary
		_, tagExists := tagMap[tag.ID]
		if !tagExists {
			tagMap[tag.ID] = ExpandedTag{
				Tag:      &tag,
				Things:   []string{},
				Channels: []string{},
				Models:   []string{},
			}
		}

		for _, appliedTag := range tag.Applications() {
			// Add the thing ID, channel ID, or model ID to the response
			if thing, ok := appliedTag.Thing(); ok {
				t := tagMap[tag.ID]
				t.Things = append(t.Things, thing.ID)
				tagMap[tag.ID] = t
			}
			if channel, ok := appliedTag.Channel(); ok {
				t := tagMap[tag.ID]
				t.Channels = append(t.Channels, channel.ID)
				tagMap[tag.ID] = t
			}
			if model, ok := appliedTag.Model(); ok {
				t := tagMap[tag.ID]
				t.Models = append(t.Models, model.ID)
				tagMap[tag.ID] = t
			}
		}

		// // Remove tags that do not apply to the current category
		// for tagID, tag := range tagMap {
		// 	if request.Input.Category == "things" && len(tag.Things) <= 0 {
		// 		delete(tagMap, tagID)
		// 	}
		// 	if request.Input.Category == "channels" && len(tag.Channels) <= 0 {
		// 		delete(tagMap, tagID)
		// 	}
		// 	if request.Input.Category == "models" && len(tag.Models) <= 0 {
		// 		delete(tagMap, tagID)
		// 	}
		// }
	}

	output := &Output{
		Tags: tagMap,
	}

	return output, nil
}
