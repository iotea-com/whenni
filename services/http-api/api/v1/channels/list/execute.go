package channelsList

import (
	"fmt"
	"strconv"

	ioteahttp "github.com/iotea-com/iotea/libs/http"
	"github.com/iotea-com/iotea/prisma/db"
	"github.com/iotea-com/iotea/services/http-api/services/prisma"
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
	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Count channels")
	var channelsCountResponse []struct {
		Count db.RawString `json:"channels_count"`
	}
	countQuery := `
		SELECT count(*) as channels_count 
		FROM app.channels 
		WHERE channels."spaceId" = $1`

	queryParams := []interface{}{request.Input.SpaceId}
	paramCount := 1

	if request.Input.Filter != "" {
		paramCount++
		countQuery += fmt.Sprintf(` AND (
			LOWER(channels.id) LIKE LOWER($%d) OR 
			LOWER(channels.name) LIKE LOWER($%d)
		)`, paramCount, paramCount)
		queryParams = append(queryParams, "%"+request.Input.Filter+"%")
	}

	if len(request.Input.TagFilter) > 0 {
		paramCount++
		countQuery += fmt.Sprintf(` AND channels.id IN (
			SELECT DISTINCT "channelId" FROM app.applied_tags
			WHERE "tagId" = ANY($%d)
		)`, paramCount)
		queryParams = append(queryParams, request.Input.TagFilter)
	}

	err := prisma.Client.Prisma.QueryRaw(countQuery, queryParams...).Exec(dbCtx, &channelsCountResponse)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error counting channels in space in the database: %s", err)),
		)
		dbSpan.End()

		return nil, err
	}

	channelsCount, err := strconv.ParseInt(string(channelsCountResponse[0].Count), 10, 16)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "parse_json"),
			attribute.String("error.message", fmt.Sprintf("could not parse count response to int - %#v: %s", channelsCountResponse, err)),
		)
		dbSpan.End()

		return nil, err
	}

	dbSpan.SetAttributes(
		attribute.Int64("channels_count", channelsCount),
	)
	dbSpan.End()

	// List channels (based on the filter)
	dbCtx, dbSpan = otel.Tracer("prisma").Start(request.Context, "List channels")
	var conditions []db.ChannelWhereParam = []db.ChannelWhereParam{
		db.Channel.SpaceID.Equals(request.Input.SpaceId),
	}

	if request.Input.Filter != "" {
		conditions = append(conditions, db.Channel.Or(
			db.Channel.ID.Mode(db.QueryModeInsensitive),
			db.Channel.Name.Mode(db.QueryModeInsensitive),
			db.Channel.ID.Contains(request.Input.Filter),
			db.Channel.Name.Contains(request.Input.Filter),
		))
	}

	if len(request.Input.TagFilter) > 0 {
		conditions = append(conditions, db.Channel.Tags.Some(db.AppliedTag.TagID.In(request.Input.TagFilter)))
	}

	channels, err := prisma.Client.Channel.FindMany(
		conditions...,
	).With(
		db.Channel.Tags.Fetch().With(
			db.AppliedTag.Tag.Fetch(),
		),
	).OrderBy(db.Channel.UpdatedAt.Order(db.DESC)).
		Take(request.Input.ResultsPerPage).
		Skip((request.Input.Page - 1) * request.Input.ResultsPerPage).
		Exec(dbCtx)
	if err != nil {
		if err.Error() == "ErrNotFound" {
			request.Span.SetAttributes(
				attribute.String("error.type", "http_request"),
				attribute.String("error.message", fmt.Sprintf("no channels found in space with ID %s", request.Input.SpaceId)),
			)
			dbSpan.End()

			output := Output{
				Channels: []db.ChannelModel{},
			}

			return &output, nil
		}

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
