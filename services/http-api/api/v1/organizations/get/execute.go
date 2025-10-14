package organizationsGet

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	"github.com/iotea-com/iotea/prisma/db"
	"github.com/iotea-com/iotea/services/http-api/services/prisma"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func execute(request *ioteahttp.Request[Input]) (*Output, error) {
	request.Span.AddEvent("execute")
	request.Span.SetAttributes(
		attribute.String("request.Input.OrgId", request.Input.OrgId),
	)

	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Get organization")
	organization, err := prisma.Client.Organization.FindUnique(
		db.Organization.ID.Equals(request.Input.OrgId),
	).With(
		db.Organization.Spaces.Fetch(),
	).Exec(dbCtx)
	if err != nil {
		if err.Error() == "ErrNotFound" {
			dbSpan.SetAttributes(
				attribute.String("error.type", "database"),
				attribute.String("error.message", fmt.Sprintf("no organization found with ID %s", request.Input.OrgId)),
			)
			dbSpan.End()
			return nil, fiber.NewError(fiber.StatusBadRequest)
		}

		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error getting organization from the database: %s", err)),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	output := Output{
		Organization: organization,
	}

	return &output, nil
}
