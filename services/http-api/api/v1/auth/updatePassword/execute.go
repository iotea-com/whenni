package updatePassword

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	"github.com/iotea-com/iotea/prisma/db"
	"github.com/iotea-com/iotea/services/http-api/services/prisma"
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
	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Update user password")
	_, err = prisma.Client.User.FindUnique(
		db.User.ID.Equals(request.Actor.Id),
	).Update(db.User.Password.Set(string(hashedPassword))).Exec(dbCtx)
	if err != nil {
		if err.Error() == "ErrNotFound" {
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
