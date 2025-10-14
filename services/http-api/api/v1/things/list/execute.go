package thingsList

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
	if request.Input.ThingCategory != nil {
		request.Span.SetAttributes(
			attribute.String("request.Input.ThingCategory", string(*request.Input.ThingCategory)),
		)
	}

	// Get the count of things (based on the filter)
	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Count things")
	var thingsCountResponse []struct {
		Count db.RawString `json:"things_count"`
	}

	countQuery := `
		SELECT count(*) as things_count 
		FROM app.things 
		WHERE things."spaceId" = $1 
		AND things.internal = false`

	queryParams := []interface{}{request.Input.SpaceId}
	paramCount := 1

	if request.Input.ThingCategory != nil {
		paramCount++
		countQuery += fmt.Sprintf(` AND things."thingCategory" = $%d`, paramCount)
		queryParams = append(queryParams, *request.Input.ThingCategory)
	}

	if request.Input.Filter != "" {
		paramCount++
		countQuery += fmt.Sprintf(` AND (
			LOWER(things.id) LIKE LOWER($%d) OR 
			LOWER(things.name) LIKE LOWER($%d) OR 
			LOWER(things."thingCategory") LIKE LOWER($%d)
		)`, paramCount, paramCount, paramCount)
		queryParams = append(queryParams, "%"+request.Input.Filter+"%")
	}

	if len(request.Input.TagFilter) > 0 {
		paramCount++
		countQuery += fmt.Sprintf(` AND things.id IN (
			SELECT DISTINCT "thingId" FROM app.applied_tags 
			WHERE "tagId" = ANY($%d)
		)`, paramCount)
		queryParams = append(queryParams, request.Input.TagFilter)
	}

	err := prisma.Client.Prisma.QueryRaw(countQuery, queryParams...).Exec(dbCtx, &thingsCountResponse)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error counting things in space in the database: %s", err)),
		)
		dbSpan.End()

		return nil, err
	}

	thingsCount, err := strconv.ParseInt(string(thingsCountResponse[0].Count), 10, 16)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "parse_json"),
			attribute.String("error.message", fmt.Sprintf("could not parse count response to int - %#v: %s", thingsCountResponse, err)),
		)
		dbSpan.End()

		return nil, err
	}

	dbSpan.SetAttributes(
		attribute.Int64("things_count", thingsCount),
	)
	dbSpan.End()

	// List the things (based on the filter)
	dbCtx, dbSpan = otel.Tracer("prisma").Start(request.Context, "List things")

	var conditions []db.ThingWhereParam = []db.ThingWhereParam{
		db.Thing.SpaceID.Equals(request.Input.SpaceId),
		db.Thing.Internal.Equals(false),
		db.Thing.ThingCategory.EqualsIfPresent(request.Input.ThingCategory.StringOrNil()),
	}

	if request.Input.Filter != "" {
		conditions = append(conditions, db.Thing.Or(
			db.Thing.ID.Mode(db.QueryModeInsensitive),
			db.Thing.Name.Mode(db.QueryModeInsensitive),
			db.Thing.ThingCategory.Mode(db.QueryModeInsensitive),
			db.Thing.ID.Contains(request.Input.Filter),
			db.Thing.Name.Contains(request.Input.Filter),
			db.Thing.ThingCategory.Contains(request.Input.Filter),
		))
	}

	if len(request.Input.TagFilter) > 0 {
		conditions = append(conditions, db.Thing.Tags.Some(db.AppliedTag.TagID.In(request.Input.TagFilter)))
	}

	things, err := prisma.Client.Thing.FindMany(
		conditions...,
	).With(
		db.Thing.Tags.Fetch().With(
			db.AppliedTag.Tag.Fetch(),
		),
	).OrderBy(db.Thing.ThingCategory.Order(db.DESC)).
		Take(request.Input.ResultsPerPage).
		Skip((request.Input.Page - 1) * request.Input.ResultsPerPage).
		Exec(dbCtx)
	if err != nil {
		if err.Error() == "ErrNotFound" {
			dbSpan.SetAttributes(
				attribute.String("error.type", "database"),
				attribute.String("error.message", fmt.Sprintf("no things found in space with ID %s", request.Input.SpaceId)),
			)
			dbSpan.End()

			output := Output{
				Things: []db.ThingModel{},
			}

			return &output, nil
		}

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
