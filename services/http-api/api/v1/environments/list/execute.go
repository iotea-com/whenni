package environmentsList

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	"github.com/iotea-com/iotea/prisma/db"
	"github.com/iotea-com/iotea/services/http-api/services/prisma"
	"go.opentelemetry.io/otel"
)

func execute(request *ioteahttp.Request[Input]) (*Output, error) {
	request.Span.AddEvent("execute")

	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Count dev environments")
	var devEnvironmentsCountResponse []struct {
		Count db.RawString `json:"environments_count"`
	}
	err := prisma.Client.Prisma.QueryRaw(`SELECT count(*) as environments_count FROM app."devEnvironments" WHERE "devEnvironments"."spaceId" = $1`, request.Input.SpaceId).Exec(dbCtx, &devEnvironmentsCountResponse)
	if err != nil {
		dbSpan.AddEvent(fmt.Sprintf("error counting dev environments in space in the database: %s", err))
		dbSpan.End()

		return nil, err
	}

	environmentsCount, err := strconv.ParseInt(string(devEnvironmentsCountResponse[0].Count), 10, 16)
	if err != nil {
		dbSpan.AddEvent(fmt.Sprintf("could not parse count response to int - %#v: %s", devEnvironmentsCountResponse, err))

		return nil, err
	}

	dbSpan.AddEvent(fmt.Sprintf("successfully counted %d dev environments in space in the database", environmentsCount))
	dbSpan.End()

	dbCtx, dbSpan = otel.Tracer("prisma").Start(request.Context, "List DevEnvironment")
	environments, err := prisma.Client.DevEnvironment.FindMany(
		db.DevEnvironment.SpaceID.Equals(request.Input.SpaceId),
	).Take(request.Input.ResultsPerPage).
		Skip((request.Input.Page - 1) * request.Input.ResultsPerPage).
		Exec(dbCtx)
	if err != nil {
		if err.Error() == "ErrNotFound" {
			dbSpan.AddEvent(fmt.Sprintf("no devenv found with spaceID %s", request.Input.SpaceId))
			dbSpan.End()
			return nil, fiber.NewError(fiber.StatusNotFound)
		}

		dbSpan.AddEvent(fmt.Sprintf("error getting environments from the database: %s", err))
		dbSpan.End()
		return nil, err
	}

	dbSpan.AddEvent("successfully got devenv from the database")
	dbSpan.End()

	output := Output{
		EnvironmentsList: environments,
		Page:             request.Input.Page,
		TotalPages:       (int(environmentsCount) + request.Input.ResultsPerPage - 1) / request.Input.ResultsPerPage,
		TotalResults:     int(environmentsCount),
		ResultsPerPage:   request.Input.ResultsPerPage,
	}

	return &output, nil
}
