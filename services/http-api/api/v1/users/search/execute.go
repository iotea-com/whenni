package search

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	"github.com/iotea-com/iotea/services/http-api/services/sqlc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func execute(request *ioteahttp.Request[Input]) (*Output, error) {
	request.Span.AddEvent("execute")
	request.Span.SetAttributes(
		attribute.String("request.Input.OrgId", request.Input.OrgId),
		attribute.String("request.Input.Query", request.Input.Query),
	)

	dbCtx, dbSpan := otel.Tracer("sqlc").Start(request.Context, "Search users")
	users, err := sqlc.Queries.SearchUsers(dbCtx, request.Input.OrgId, &request.Input.Query)
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error searching users in the database: %s", err)),
		)
		dbSpan.End()
		return nil, err
	}
	if len(users) == 0 {
		request.Span.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("no users found that match query %s", request.Input.Query)),
		)
		dbSpan.End()
		return nil, fiber.NewError(fiber.StatusBadRequest)
	}

	dbSpan.End()

	output := Output{
		Users: users,
	}

	return &output, nil
}
