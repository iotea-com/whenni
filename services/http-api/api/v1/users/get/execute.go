package get

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	sqldb "github.com/iotea-com/iotea/db/sqlc"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	"github.com/iotea-com/iotea/services/http-api/services/sqlc"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func execute(request *ioteahttp.Request[Input]) (*Output, error) {
	request.Span.AddEvent("execute")
	request.Span.SetAttributes(
		attribute.String("request.Input.UserId", request.Input.UserId),
	)

	dbCtx, dbSpan := otel.Tracer("sqlc").Start(request.Context, "Get profile")
	user, err := sqlc.Queries.GetUser(dbCtx, request.Input.UserId)
	if err != nil {
		if err == pgx.ErrNoRows {
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

	members, err := sqlc.Queries.ListUserOrganizationMembers(dbCtx, request.Input.UserId)
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error getting user organizations from the database: %s", err)),
		)
		dbSpan.End()
		return nil, err
	}

	spacesByOrg := make(map[string][]sqldb.AppSpace, len(members))
	for _, member := range members {
		spaces, listErr := sqlc.Queries.ListSpaces(dbCtx, member.OrganizationID)
		if listErr != nil {
			request.Span.SetAttributes(
				attribute.String("error.type", "database"),
				attribute.String("error.message", fmt.Sprintf("error getting spaces for organization %s: %s", member.OrganizationID, listErr)),
			)
			dbSpan.End()
			return nil, listErr
		}
		spacesByOrg[member.OrganizationID] = spaces
	}

	dbSpan.End()

	output := Output{
		User: toUser(user, members, spacesByOrg),
	}

	return &output, nil
}
