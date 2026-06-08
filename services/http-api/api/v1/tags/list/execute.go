package tagsList

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
		attribute.String("request.Input.SpaceId", request.Input.SpaceId),
		attribute.String("request.Input.Category", request.Input.Category),
	)

	// List tags in space (based on the filter)
	dbCtx, dbSpan := otel.Tracer("sqlc").Start(request.Context, "List tags in space")

	tags, err := sqlc.Queries.ListTags(dbCtx, request.Input.SpaceId)
	if err != nil {
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	dbSpan.SetAttributes(
		attribute.String("response.Tags", fmt.Sprintf("%v", tags)),
	)

	// Group applied tags by tag
	tagMap := make(map[string]ExpandedTag, len(tags))
	for _, tag := range tags {
		expandedTag := ExpandedTag{
			Tag:      tag,
			Things:   []string{},
			Channels: []string{},
			Models:   []string{},
		}

		appliedTags, listAppliedTagsErr := sqlc.Queries.ListAppliedTagsByTag(dbCtx, tag.ID)
		if listAppliedTagsErr != nil {
			return nil, listAppliedTagsErr
		}
		for _, appliedTag := range appliedTags {
			if appliedTag.ThingID != nil {
				expandedTag.Things = append(expandedTag.Things, *appliedTag.ThingID)
			}
			if appliedTag.ChannelID != nil {
				expandedTag.Channels = append(expandedTag.Channels, *appliedTag.ChannelID)
			}
			if appliedTag.ModelID != nil {
				expandedTag.Models = append(expandedTag.Models, *appliedTag.ModelID)
			}
		}
		tagMap[tag.ID] = expandedTag
	}

	output := &Output{
		Tags: tagMap,
	}

	return output, nil
}
