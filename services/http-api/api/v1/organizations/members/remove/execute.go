package organizationsMembersRemove

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
		attribute.String("request.Input.UserId", request.Input.UserId),
	)

	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Delete member from space")
	_, err := prisma.Client.OrganizationMember.FindUnique(
		db.OrganizationMember.OrganizationIDUserID(
			db.OrganizationMember.OrganizationID.Equals(request.Input.OrgId),
			db.OrganizationMember.UserID.Equals(request.Input.UserId),
		),
	).Delete().Exec(dbCtx)

	if err != nil {
		errMessage := fmt.Sprintf("error deleting member in organization with ID %s from the database: %s", request.Input.OrgId, err)
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", errMessage),
		)
		dbSpan.End()

		errors := []any{errMessage}
		response := ioteahttp.NewErrorResponse(errors)
		responseJson, err := response.MarshalJson()
		if err != nil {
			return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		return nil, fiber.NewError(fiber.StatusBadRequest, string(responseJson))
	}

	dbSpan.End()

	return nil, nil
}
