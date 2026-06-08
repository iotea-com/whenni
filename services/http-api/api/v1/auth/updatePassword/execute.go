package updatePassword

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	"github.com/iotea-com/iotea/services/http-api/services/sqlc"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"golang.org/x/crypto/bcrypt"
)

func execute(request *ioteahttp.Request[Input]) (*Output, error) {
	request.Span.AddEvent("execute")

	// Bcrypt the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "bcrypt"),
			attribute.String("error.message", fmt.Sprintf("error hashing password: %s", err)),
		)
		errResponse := ioteahttp.NewErrorResponse([]any{"Unable to change password. Please try again."})
		marshalledErrResponse, _ := errResponse.MarshalJson()
		return nil, fiber.NewError(fiber.StatusInternalServerError, string(marshalledErrResponse))
	}

	request.Span.SetAttributes(
		attribute.String("hashedPassword", string(hashedPassword)),
		attribute.String("userId", request.Actor.Id),
	)

	// Update the user's password
	dbCtx, dbSpan := otel.Tracer("sqlc").Start(request.Context, "Update user password")
	updatedPassword := string(hashedPassword)
	_, err = sqlc.Queries.UpdateUserPassword(dbCtx, request.Actor.Id, &updatedPassword)
	if err != nil {
		if err == pgx.ErrNoRows {
			dbSpan.SetAttributes(
				attribute.String("error.type", "database"),
				attribute.String("error.message", fmt.Sprintf("error updating user password: %s", err)),
			)
			dbSpan.End()
			errResponse := ioteahttp.NewErrorResponse([]any{"Invalid magic link"})
			marshalledErrResponse, _ := errResponse.MarshalJson()
			return nil, fiber.NewError(fiber.StatusUnauthorized, string(marshalledErrResponse))
		}

		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error deleting API key from the database: %s", err)),
		)
		dbSpan.End()
		errResponse := ioteahttp.NewErrorResponse([]any{"Unable to verify magic link. Please try again."})
		marshalledErrResponse, _ := errResponse.MarshalJson()
		return nil, fiber.NewError(fiber.StatusInternalServerError, string(marshalledErrResponse))
	}

	dbSpan.End()

	return &Output{}, nil
}
