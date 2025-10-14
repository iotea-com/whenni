package get

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
		attribute.String("request.Input.UserId", request.Input.UserId),
	)

	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Get profile")
	user, err := prisma.Client.User.FindUnique(
		db.User.ID.Equals(request.Input.UserId),
	).Omit(
		db.User.Password.Field(),
	).With(
		db.User.Organizations.Fetch().With(
			db.OrganizationMember.Organization.Fetch().With(
				db.Organization.Spaces.Fetch(),
			),
		),
	).Exec(dbCtx)

	if err != nil {
		if err.Error() == "ErrNotFound" {
			request.Span.SetAttributes(
				attribute.String("error.type", "database"),
				attribute.String("error.message", fmt.Sprintf("no user found with ID %s", request.Input.UserId)),
			)
			dbSpan.End()
			return nil, fiber.NewError(fiber.StatusBadRequest)
		}

		request.Span.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error getting user from the database: %s", err)),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	output := Output{
		User: user,
	}

	return &output, nil
}
