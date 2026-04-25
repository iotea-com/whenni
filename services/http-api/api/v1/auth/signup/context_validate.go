package signup

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	"github.com/iotea-com/iotea/services/http-api/services/sqlc"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/attribute"
)

func contextValidate(request *ioteahttp.Request[Input]) error {
	request.Span.AddEvent("contextValidate")

	// Ensure the email is not already in use
	userEmail := request.Input.Email
	_, err := sqlc.Queries.GetUserByEmail(request.Context, &userEmail)
	if err != nil && err != pgx.ErrNoRows {
		request.Span.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("could not validate email: %s", err)),
		)
		errResponse := ioteahttp.NewErrorResponse([]any{"Could not validate email. Please try again."})
		marshalledErrResponse, _ := errResponse.MarshalJson()
		return fiber.NewError(fiber.StatusInternalServerError, string(marshalledErrResponse))
	}

	if err == nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "email_already_in_use"),
			attribute.String("error.message", fmt.Sprintf("email '%s' already in use", request.Input.Email)),
		)
		errMessage := ioteahttp.NewErrorResponse([]any{"Email is already in use"})
		marshalledErrResponse, _ := errMessage.MarshalJson()
		return fiber.NewError(fiber.StatusConflict, string(marshalledErrResponse))
	}

	request.Span.SetAttributes(
		attribute.String("contextValidate.result", "pass"),
	)
	return nil
}
