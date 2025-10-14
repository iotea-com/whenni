package magicLinkVerify

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	ioteahttputil "github.com/iotea-com/iotea/libs/http/util"
	"github.com/iotea-com/iotea/prisma/db"
	"github.com/iotea-com/iotea/services/http-api/config"
	"github.com/iotea-com/iotea/services/http-api/services/prisma"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func execute(request *ioteahttp.Request[Input]) (*Output, error) {
	request.Span.AddEvent("execute")
	request.Span.SetAttributes(
		attribute.String("request.Input.VerificationToken", request.Input.VerificationToken),
	)

	// Retrieve the verification token from the database
	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Retrieve verification token")
	verificationToken, err := prisma.Client.AuthToken.FindUnique(
		db.AuthToken.Token.Equals(request.Input.VerificationToken),
	).With(db.AuthToken.User.Fetch()).Exec(dbCtx)
	if err != nil {
		if err.Error() == "ErrNotFound" {
			dbSpan.SetAttributes(
				attribute.String("error.type", "verification_token_not_found"),
				attribute.String("error.message", fmt.Sprintf("verification token not found: %s", request.Input.VerificationToken)),
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

	// Check if the verification token has expired
	if verificationToken.ExpiresAt.Before(time.Now()) {
		request.Span.SetAttributes(
			attribute.String("error.type", "verification_token_expired"),
			attribute.String("error.message", fmt.Sprintf("verification token expired: %s", verificationToken.Token)),
		)
		errResponse := ioteahttp.NewErrorResponse([]any{"Magic link has expired. Please request a new one."})
		marshalledErrResponse, _ := errResponse.MarshalJson()
		return nil, fiber.NewError(fiber.StatusUnauthorized, string(marshalledErrResponse))
	}

	// Verify user email if not verified
	_, isEmailVerified := verificationToken.User().EmailVerified()
	if !isEmailVerified {
		dbCtx, dbSpan = otel.Tracer("prisma").Start(request.Context, "Set user email as verified")
		_, err := prisma.Client.User.FindUnique(
			db.User.ID.Equals(verificationToken.User().ID),
		).Update(db.User.EmailVerified.Set(time.Now().UTC())).Exec(dbCtx)
		if err != nil {
			dbSpan.SetAttributes(
				attribute.String("error.type", "database"),
				attribute.String("error.message", fmt.Sprintf("error updating user email as verified: %s", err)),
			)
			dbSpan.End()
			errResponse := ioteahttp.NewErrorResponse([]any{"Unable to access the database. Please try again."})
			marshalledErrResponse, _ := errResponse.MarshalJson()
			return nil, fiber.NewError(fiber.StatusInternalServerError, string(marshalledErrResponse))
		}

		dbSpan.End()
	}

	// Create a new access token
	accessToken, err := ioteahttputil.GenerateAccessToken(verificationToken.User().ID, config.VaultConf.JwtSecret)
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "jwt_token_generation"),
			attribute.String("error.message", fmt.Sprintf("error generating JWT token: %s", err)),
		)
		errResponse := ioteahttp.NewErrorResponse([]any{"Unable to create a session token. Please try again."})
		marshalledErrResponse, _ := errResponse.MarshalJson()
		return nil, fiber.NewError(fiber.StatusInternalServerError, string(marshalledErrResponse))
	}

	// Create a new refresh token
	refreshToken, err := ioteahttputil.GenerateRefreshToken()
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "refresh_token_generation"),
			attribute.String("error.message", fmt.Sprintf("error generating refresh token: %s", err)),
		)
		errResponse := ioteahttp.NewErrorResponse([]any{"Unable to create a session token. Please try again."})
		marshalledErrResponse, _ := errResponse.MarshalJson()
		return nil, fiber.NewError(fiber.StatusInternalServerError, string(marshalledErrResponse))
	}

	output := &Output{
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
	}

	// Clear old refresh tokens from the database
	dbCtx, dbSpan = otel.Tracer("prisma").Start(request.Context, "Clear old refresh tokens")
	_, err = prisma.Client.AuthToken.FindMany(
		db.AuthToken.UserID.Equals(verificationToken.User().ID),
		db.AuthToken.Type.Equals(db.AuthTokenTypeRefresh),
	).Delete().Exec(dbCtx)

	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error clearing old refresh tokens in the database: %s", err)),
		)
		dbSpan.End()
		errMessage := ioteahttp.NewErrorResponse([]any{"Could not store refresh token"})
		marshalledErrMessage, _ := errMessage.MarshalJson()
		return nil, fiber.NewError(fiber.StatusInternalServerError, string(marshalledErrMessage))
	}

	dbSpan.End()

	// Store the refresh token in the database
	dbCtx, dbSpan = otel.Tracer("prisma").Start(request.Context, "Store refresh token")
	_, err = prisma.Client.AuthToken.CreateOne(
		db.AuthToken.Token.Set(*refreshToken),
		db.AuthToken.Type.Set(db.AuthTokenTypeRefresh),
		db.AuthToken.ExpiresAt.Set(time.Now().Add(ioteahttputil.RefreshTokenExpiry).UTC()),
		db.AuthToken.User.Link(db.User.ID.Equals(verificationToken.User().ID)),
	).Exec(dbCtx)

	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error storing refresh token: %s", err)),
		)
		dbSpan.End()
		errResponse := ioteahttp.NewErrorResponse([]any{"Unable to verify magic link. Please try again."})
		marshalledErrResponse, _ := errResponse.MarshalJson()
		return nil, fiber.NewError(fiber.StatusInternalServerError, string(marshalledErrResponse))
	}

	dbSpan.End()

	// Remove all verification tokens for the user from the database
	dbCtx, dbSpan = otel.Tracer("prisma").Start(request.Context, "Remove verification token")
	_, err = prisma.Client.AuthToken.FindMany(
		db.AuthToken.User.Where(db.User.ID.Equals(verificationToken.User().ID)),
		db.AuthToken.Type.Equals(db.AuthTokenTypeMagicLink),
	).Delete().Exec(dbCtx)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error removing verification token from the database: %s", err)),
		)
		dbSpan.End()
		errResponse := ioteahttp.NewErrorResponse([]any{"Unable to verify magic link. Please try again."})
		marshalledErrResponse, _ := errResponse.MarshalJson()
		return nil, fiber.NewError(fiber.StatusInternalServerError, string(marshalledErrResponse))
	}

	dbSpan.End()

	return output, nil
}
