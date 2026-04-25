package thingsList

import (
	"fmt"

	sqldb "github.com/iotea-com/iotea/db/sqlc"
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
	if request.Input.ThingCategory != nil {
		request.Span.SetAttributes(
			attribute.String("request.Input.ThingCategory", string(*request.Input.ThingCategory)),
		)
	}

	// Get the count of things (based on the filter)
	dbCtx, dbSpan := otel.Tracer("sqlc").Start(request.Context, "Count things")
	thingCategory := ""
	if request.Input.ThingCategory != nil {
		thingCategory = request.Input.ThingCategory.String()
	}
	filter := ""
	if request.Input.Filter != "" {
		filter = "%" + request.Input.Filter + "%"
	}
	var tagFilter []string
	if len(request.Input.TagFilter) > 0 {
		tagFilter = request.Input.TagFilter
	}

	thingsCount, err := sqlc.Queries.CountThingsFiltered(dbCtx, request.Input.SpaceId, thingCategory, filter, tagFilter)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error counting things in space in the database: %s", err)),
		)
		dbSpan.End()

		return nil, err
	}

	dbSpan.SetAttributes(
		attribute.Int64("things_count", thingsCount),
	)
	dbSpan.End()

	// List the things (based on the filter)
	dbCtx, dbSpan = otel.Tracer("sqlc").Start(request.Context, "List things")
	things, err := sqlc.Queries.ListThingsFiltered(dbCtx, sqldb.ListThingsFilteredParams{
		SpaceID: request.Input.SpaceId,
		Column2: thingCategory,
		Column3: filter,
		Column4: tagFilter,
		Limit:   int32(request.Input.ResultsPerPage),
		Offset:  int32((request.Input.Page - 1) * request.Input.ResultsPerPage),
	})
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error listing things in the database: %s", err)),
		)
		dbSpan.End()

		return nil, err
	}

	dbSpan.End()

	output := Output{
		Things:         things,
		Page:           request.Input.Page,
		TotalPages:     (int(thingsCount) + request.Input.ResultsPerPage - 1) / request.Input.ResultsPerPage,
		TotalResults:   int(thingsCount),
		ResultsPerPage: request.Input.ResultsPerPage,
	}

	return &output, nil
}
