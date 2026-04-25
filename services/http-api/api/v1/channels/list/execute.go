package channelsList

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
		attribute.Int("request.Input.Page", request.Input.Page),
		attribute.Int("request.Input.ResultsPerPage", request.Input.ResultsPerPage),
		attribute.String("request.Input.SpaceId", request.Input.SpaceId),
		attribute.String("request.Input.Filter", request.Input.Filter),
		attribute.String("request.Input.TagFilter", fmt.Sprintf("%#v", request.Input.TagFilter)),
	)

	// Get the count of channels (based on the filter)
	dbCtx, dbSpan := otel.Tracer("sqlc").Start(request.Context, "Count channels")
	channelsCount, err := sqlc.Queries.CountChannelsFiltered(
		dbCtx,
		request.Input.SpaceId,
		request.Input.Filter,
		request.Input.TagFilter,
	)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error counting channels in space in the database: %s", err)),
		)
		dbSpan.End()

		return nil, err
	}

	dbSpan.SetAttributes(
		attribute.Int64("channels_count", channelsCount),
	)
	dbSpan.End()

	// List channels (based on the filter)
	dbCtx, dbSpan = otel.Tracer("sqlc").Start(request.Context, "List channels")
	channels, err := sqlc.Queries.ListChannelsFiltered(
		dbCtx,
		request.Input.SpaceId,
		request.Input.Filter,
		request.Input.TagFilter,
		int32(request.Input.ResultsPerPage),
		int32((request.Input.Page-1)*request.Input.ResultsPerPage),
	)
	if err != nil {
		dbSpan.End()

		return nil, err
	}

	dbSpan.End()

	output := Output{
		Channels:       channels,
		Page:           request.Input.Page,
		TotalPages:     (int(channelsCount) + request.Input.ResultsPerPage - 1) / request.Input.ResultsPerPage,
		TotalResults:   int(channelsCount),
		ResultsPerPage: request.Input.ResultsPerPage,
	}

	return &output, nil
}
