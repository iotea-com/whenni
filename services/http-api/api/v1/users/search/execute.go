package search

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
		attribute.String("request.Input.Query", request.Input.Query),
	)

	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Search users")
	users, err := prisma.Client.User.FindMany(
		db.User.Email.Contains(request.Input.Query),
		db.User.Organizations.None(
			db.OrganizationMember.OrganizationID.Equals(request.Input.OrgId),
		),
	).OrderBy(
		db.User.Relevance_.Fields([]db.UserOrderByRelevanceFieldEnum{
			db.UserOrderByRelevanceFieldEnumEmail,
		}),
		db.User.Relevance_.Search("database"),
		db.User.Relevance_.Sort("desc"),
	).Exec(dbCtx)
	if err != nil {
		if err.Error() == "ErrNotFound" {
			request.Span.SetAttributes(
				attribute.String("error.type", "database"),
				attribute.String("error.message", fmt.Sprintf("no users found that match query %s", request.Input.Query)),
			)
			dbSpan.End()
			return nil, fiber.NewError(fiber.StatusBadRequest)
		}

		request.Span.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error searching users in the database: %s", err)),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	output := Output{
		Users: users,
	}

	return &output, nil
}
