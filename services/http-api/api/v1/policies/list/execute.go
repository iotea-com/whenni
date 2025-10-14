package policiesList

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
		attribute.String("request.Input.SpaceId", request.Input.SpaceId),
		attribute.Int("request.Input.Page", request.Input.Page),
		attribute.Int("request.Input.ResultsPerPage", request.Input.ResultsPerPage),
	)

	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Count policies")
	var policiesCountResponse []struct {
		Count db.RawString `json:"policies_count"`
	}
	err := prisma.Client.Prisma.QueryRaw(`SELECT count(*) as policies_count FROM app.certificates WHERE certificates."spaceId" = $1`, request.Input.SpaceId).Exec(dbCtx, &policiesCountResponse)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error counting policies in space in the database: %s", err)),
		)
		dbSpan.End()

		return nil, err
	}

	policiesCount, err := strconv.ParseInt(string(policiesCountResponse[0].Count), 10, 16)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("could not parse count response to int - %#v: %s", policiesCountResponse, err)),
		)
		dbSpan.End()

		return nil, err
	}

	dbSpan.SetAttributes(
		attribute.Int64("policiesCount", policiesCount),
	)
	dbSpan.End()

	dbCtx, dbSpan = otel.Tracer("prisma").Start(request.Context, "List policies in the space")
	certificates, err := prisma.Client.Certificate.FindMany(
		db.Certificate.SpaceID.Equals(request.Input.SpaceId),
	).Take(request.Input.ResultsPerPage).
		Skip((request.Input.Page - 1) * request.Input.ResultsPerPage).
		Exec(dbCtx)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error listing permission sets in the database: %s", err)),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	output := Output{
		Certificates:   certificates,
		Page:           request.Input.Page,
		TotalPages:     (int(policiesCount) + request.Input.ResultsPerPage - 1) / request.Input.ResultsPerPage,
		TotalResults:   int(policiesCount),
		ResultsPerPage: request.Input.ResultsPerPage,
	}

	return &output, nil
}
