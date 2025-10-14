package modelsList

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
	)

	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Count models")
	var modelsCountResponse []struct {
		Count db.RawString `json:"models_count"`
	}
	countQuery := `
		SELECT count(*) as models_count 
		FROM app.models 
		WHERE models."spaceId" = $1`

	queryParams := []interface{}{request.Input.SpaceId}
	paramCount := 1

	if request.Input.Filter != "" {
		paramCount++
		countQuery += fmt.Sprintf(` AND (
			LOWER(models.id) LIKE LOWER($%d) OR 
			LOWER(models.name) LIKE LOWER($%d)
		)`, paramCount, paramCount)
		queryParams = append(queryParams, "%"+request.Input.Filter+"%")
	}

	if len(request.Input.TagFilter) > 0 {
		paramCount++
		countQuery += fmt.Sprintf(` AND models.id IN (
			SELECT DISTINCT "modelId" FROM app.applied_tags
			WHERE "tagId" = ANY($%d)
		)`, paramCount)
		queryParams = append(queryParams, request.Input.TagFilter)
	}

	err := prisma.Client.Prisma.QueryRaw(countQuery, queryParams...).Exec(dbCtx, &modelsCountResponse)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error counting models in space in the database: %s", err)),
		)
		dbSpan.End()

		return nil, err
	}

	modelsCount, err := strconv.ParseInt(string(modelsCountResponse[0].Count), 10, 16)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "parse_json"),
			attribute.String("error.message", fmt.Sprintf("could not parse count response to int - %#v: %s", modelsCountResponse, err)),
		)
		dbSpan.End()

		return nil, err
	}

	dbSpan.SetAttributes(
		attribute.Int64("models_count", modelsCount),
	)
	dbSpan.End()

	dbCtx, dbSpan = otel.Tracer("prisma").Start(request.Context, "List models")
	var conditions []db.ModelWhereParam = []db.ModelWhereParam{
		db.Model.SpaceID.Equals(request.Input.SpaceId),
	}

	if request.Input.Filter != "" {
		conditions = append(conditions, db.Model.Or(
			db.Model.ID.Mode(db.QueryModeInsensitive),
			db.Model.Name.Mode(db.QueryModeInsensitive),
			db.Model.ID.Contains(request.Input.Filter),
			db.Model.Name.Contains(request.Input.Filter),
		))
	}

	if len(request.Input.TagFilter) > 0 {
		conditions = append(conditions, db.Model.Tags.Some(db.AppliedTag.TagID.In(request.Input.TagFilter)))
	}

	models, err := prisma.Client.Model.FindMany(
		conditions...,
	).With(
		db.Model.Tags.Fetch().With(
			db.AppliedTag.Tag.Fetch(),
		),
	).OrderBy(db.Model.UpdatedAt.Order(db.DESC)).
		Take(request.Input.ResultsPerPage).
		Skip((request.Input.Page - 1) * request.Input.ResultsPerPage).
		Exec(dbCtx)
	if err != nil {
		if err.Error() == "ErrNotFound" {
			dbSpan.SetAttributes(
				attribute.String("error.type", "database"),
				attribute.String("error.message", fmt.Sprintf("no models found in space with ID %s", request.Input.SpaceId)),
			)
			dbSpan.End()

			output := &Output{
				Models: []db.ModelModel{},
			}

			return output, nil
		}

		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error listing models in the database: %s", err)),
		)
		dbSpan.End()

		return nil, err
	}

	dbSpan.End()

	output := &Output{
		Models:         models,
		Page:           request.Input.Page,
		TotalPages:     (int(modelsCount) + request.Input.ResultsPerPage - 1) / request.Input.ResultsPerPage,
		TotalResults:   int(modelsCount),
		ResultsPerPage: request.Input.ResultsPerPage,
	}

	return output, nil
}
