package refresh

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
		attribute.String("request.Input.RefreshToken", request.Input.RefreshToken),
	)

	// Get the refresh token and user from the database
	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Get refresh token with user")
	refreshToken, err := prisma.Client.AuthToken.FindUnique(
		db.AuthToken.Token.Equals(request.Input.RefreshToken),
	).With(db.AuthToken.User.Fetch()).Exec(dbCtx)

	if err != nil {
		if err.Error() == "ErrNotFound" {
			dbSpan.SetAttributes(
				attribute.String("error.type", "refresh_token_not_found"),
				attribute.String("error.message", fmt.Sprintf("refresh token not found: %s", request.Input.RefreshToken)),
			)
			dbSpan.End()
			errMessage := ioteahttp.NewErrorResponse([]any{"Invalid refresh token"})
			marshalledErrMessage, _ := errMessage.MarshalJson()
			return nil, fiber.NewError(fiber.StatusPreconditionFailed, string(marshalledErrMessage))
		}

		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error getting refresh token with user from the database: %s", err)),
		)
		dbSpan.End()
		errMessage := ioteahttp.NewErrorResponse([]any{"Could not check refresh token"})
		marshalledErrMessage, _ := errMessage.MarshalJson()
		return nil, fiber.NewError(fiber.StatusInternalServerError, string(marshalledErrMessage))
	}

	dbSpan.End()

	// Check if the refresh token is expired
	if refreshToken.ExpiresAt.Before(time.Now().UTC()) {
		request.Span.SetAttributes(
			attribute.String("error.type", "refresh_token_expired"),
			attribute.String("error.message", fmt.Sprintf("refresh token expired: %s", refreshToken.ExpiresAt)),
		)
		errMessage := ioteahttp.NewErrorResponse([]any{"Refresh token expired"})
		marshalledErrMessage, _ := errMessage.MarshalJson()
		return nil, fiber.NewError(fiber.StatusUnauthorized, string(marshalledErrMessage))
	}

	// Generate a new access token
	user := refreshToken.User()
	accessToken, err := ioteahttputil.GenerateAccessToken(user.ID, config.VaultConf.JwtSecret)
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "jwt_generation"),
			attribute.String("error.message", fmt.Sprintf("error generating jwt: %s", err)),
		)
		errMessage := ioteahttp.NewErrorResponse([]any{"Could not generate access token"})
		marshalledErrMessage, _ := errMessage.MarshalJson()
		return nil, fiber.NewError(fiber.StatusInternalServerError, string(marshalledErrMessage))
	}

	// Generate a new refresh token
	newRefreshToken, err := ioteahttputil.GenerateRefreshToken()
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "refresh_token_generation"),
			attribute.String("error.message", fmt.Sprintf("error generating refresh token: %s", err)),
		)
		errMessage := ioteahttp.NewErrorResponse([]any{"Could not generate refresh token"})
		marshalledErrMessage, _ := errMessage.MarshalJson()
		return nil, fiber.NewError(fiber.StatusInternalServerError, string(marshalledErrMessage))
	}

	// Clear old refresh tokens from the database
	dbCtx, dbSpan = otel.Tracer("prisma").Start(request.Context, "Clear old refresh tokens")
	_, err = prisma.Client.AuthToken.FindMany(
		db.AuthToken.UserID.Equals(refreshToken.User().ID),
		db.AuthToken.Type.Equals(db.AuthTokenTypeRefresh),
	).Delete().Exec(dbCtx)

	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error clearing old refresh tokens in the database: %s", err)),
		)
		dbSpan.End()
		errMessage := ioteahttp.NewErrorResponse([]any{"Could not store new refresh token"})
		marshalledErrMessage, _ := errMessage.MarshalJson()
		return nil, fiber.NewError(fiber.StatusInternalServerError, string(marshalledErrMessage))
	}

	// Store the new refresh token in the database
	dbCtx, dbSpan = otel.Tracer("prisma").Start(request.Context, "Store new refresh token")
	_, err = prisma.Client.AuthToken.CreateOne(
		db.AuthToken.Token.Set(*newRefreshToken),
		db.AuthToken.Type.Set(db.AuthTokenTypeRefresh),
		db.AuthToken.ExpiresAt.Set(time.Now().Add(ioteahttputil.RefreshTokenExpiry).UTC()),
		db.AuthToken.User.Link(db.User.ID.Equals(refreshToken.User().ID)),
	).Exec(dbCtx)

	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error storing new refresh token in the database: %s", err)),
		)
		dbSpan.End()
		errMessage := ioteahttp.NewErrorResponse([]any{"Could not store new refresh token"})
		marshalledErrMessage, _ := errMessage.MarshalJson()
		return nil, fiber.NewError(fiber.StatusInternalServerError, string(marshalledErrMessage))
	}

	dbSpan.End()

	output := &Output{
		AccessToken:  *accessToken,
		RefreshToken: *newRefreshToken,
	}

	return output, nil
}
