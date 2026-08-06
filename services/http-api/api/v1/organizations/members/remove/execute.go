package organizationsMembersRemove

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	gruenthttp "github.com/ongruent/gruent/libs/http"
	"github.com/ongruent/gruent/services/http-api/services/sqlc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func execute(request *gruenthttp.Request[Input]) (*Output, error) {
	request.Span.AddEvent("execute")
	request.Span.SetAttributes(
		attribute.String("request.Input.OrgId", request.Input.OrgId),
		attribute.String("request.Input.UserId", request.Input.UserId),
	)

	dbCtx, dbSpan := otel.Tracer("sqlc").Start(request.Context, "Delete member from organization")
	member, err := sqlc.Queries.GetOrganizationMember(dbCtx, request.Input.OrgId, request.Input.UserId)
	if err == nil {
		err = sqlc.Queries.DeleteOrganizationMember(dbCtx, member.ID)
	}

	if err != nil {
		errMessage := fmt.Sprintf("error deleting member in organization with ID %s from the database: %s", request.Input.OrgId, err)
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", errMessage),
		)
		dbSpan.End()

		errors := []any{errMessage}
		response := gruenthttp.NewErrorResponse(errors)
		responseJson, err := response.MarshalJson()
		if err != nil {
			return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		return nil, fiber.NewError(fiber.StatusBadRequest, string(responseJson))
	}

	dbSpan.End()

	return nil, nil
}
