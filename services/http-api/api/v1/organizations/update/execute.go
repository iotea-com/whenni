package organizationsUpdate

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	"github.com/iotea-com/iotea/prisma/db"
	"github.com/iotea-com/iotea/services/http-api/services/prisma"
	"github.com/iotea-com/iotea/services/http-api/util"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func execute(request *ioteahttp.Request[Input]) (*Output, error) {
	request.Span.AddEvent("execute")
	request.Span.SetAttributes(
		attribute.String("request.Input.Organization.ID", request.Input.Organization.ID),
	)

	now := util.GetCurrentTime()

	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Update organization")
	organization, err := prisma.Client.Organization.FindUnique(
		db.Organization.ID.Equals(request.Input.Organization.ID),
	).Update(
		db.Organization.Name.Set(request.Input.Organization.Name),
		db.Organization.UpdatedBy.Set(request.GetActorId()),
		db.Organization.UpdatedAt.Set(now),
	).Exec(dbCtx)
	if err != nil {
		if err.Error() == "ErrNotFound" {
			dbSpan.SetAttributes(
				attribute.String("error.type", "database"),
				attribute.String("error.message", fmt.Sprintf("no organization found with ID %s", request.Input.Organization.ID)),
			)
			dbSpan.End()

			return nil, fiber.NewError(fiber.StatusBadRequest)
		}

		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error updating organization in the database: %s", err)),
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
