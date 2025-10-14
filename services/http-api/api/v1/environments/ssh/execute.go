package environmentsSsh

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	"github.com/iotea-com/iotea/prisma/db"
	"github.com/iotea-com/iotea/services/http-api/services/prisma"
	"go.opentelemetry.io/otel"
)

func execute(request *ioteahttp.Request[Input]) (*Output, error) {
	request.Span.AddEvent("execute")

	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Get DevEnvironment DNS")
	devenv, err := prisma.Client.DevEnvironment.FindUnique(
		db.DevEnvironment.EnvironmentID.Equals(request.Input.EnvironmentID),
	).Exec(dbCtx)
	if err != nil {
		if err.Error() == "ErrNotFound" {
			dbSpan.AddEvent(fmt.Sprintf("no devenv found with ID %s", request.Input.SpaceId))
			dbSpan.End()
			return nil, fiber.NewError(fiber.StatusNotFound)
		}

		dbSpan.AddEvent(fmt.Sprintf("error getting devenv from the database: %s", err))
		dbSpan.End()
		return nil, err
	}

	dbSpan.AddEvent("successfully got devenv from the database")
	dbSpan.End()

	output := Output{
		PublicDNS: devenv.PublicDNS,
	}

	return &output, nil
}
