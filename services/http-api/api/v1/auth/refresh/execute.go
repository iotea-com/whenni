package refresh

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	sqldb "github.com/ongruent/gruent/db/sqlc"
	gruenthttp "github.com/ongruent/gruent/libs/http"
	gruenthttputil "github.com/ongruent/gruent/libs/http/util"
	"github.com/ongruent/gruent/services/http-api/config"
	"github.com/ongruent/gruent/services/http-api/services/sqlc"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func execute(request *gruenthttp.Request[Input]) (*Output, error) {
	request.Span.AddEvent("execute")
	request.Span.SetAttributes(
		attribute.String("request.Input.RefreshToken", request.Input.RefreshToken),
	)

	// Get the refresh token and user from the database
	dbCtx, dbSpan := otel.Tracer("sqlc").Start(request.Context, "Get refresh token with user")
	refreshToken, err := sqlc.Queries.GetAuthTokenWithUser(dbCtx, request.Input.RefreshToken)

	if err != nil {
		if err == pgx.ErrNoRows {
			dbSpan.SetAttributes(
				attribute.String("error.type", "refresh_token_not_found"),
				attribute.String("error.message", fmt.Sprintf("refresh token not found: %s", request.Input.RefreshToken)),
			)
			dbSpan.End()
			errMessage := gruenthttp.NewErrorResponse([]any{"Invalid refresh token"})
			marshalledErrMessage, _ := errMessage.MarshalJson()
			return nil, fiber.NewError(fiber.StatusPreconditionFailed, string(marshalledErrMessage))
		}

		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error getting refresh token with user from the database: %s", err)),
		)
		dbSpan.End()
		errMessage := gruenthttp.NewErrorResponse([]any{"Could not check refresh token"})
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
		errMessage := gruenthttp.NewErrorResponse([]any{"Refresh token expired"})
		marshalledErrMessage, _ := errMessage.MarshalJson()
		return nil, fiber.NewError(fiber.StatusUnauthorized, string(marshalledErrMessage))
	}

	// Generate a new access token
	accessToken, err := gruenthttputil.GenerateAccessToken(refreshToken.JoinedUserID, config.VaultConf.JwtSecret)
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "jwt_generation"),
			attribute.String("error.message", fmt.Sprintf("error generating jwt: %s", err)),
		)
		errMessage := gruenthttp.NewErrorResponse([]any{"Could not generate access token"})
		marshalledErrMessage, _ := errMessage.MarshalJson()
		return nil, fiber.NewError(fiber.StatusInternalServerError, string(marshalledErrMessage))
	}

	// Generate a new refresh token
	newRefreshToken, err := gruenthttputil.GenerateRefreshToken()
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "refresh_token_generation"),
			attribute.String("error.message", fmt.Sprintf("error generating refresh token: %s", err)),
		)
		errMessage := gruenthttp.NewErrorResponse([]any{"Could not generate refresh token"})
		marshalledErrMessage, _ := errMessage.MarshalJson()
		return nil, fiber.NewError(fiber.StatusInternalServerError, string(marshalledErrMessage))
	}

	// Clear old refresh tokens from the database
	dbCtx, dbSpan = otel.Tracer("sqlc").Start(request.Context, "Clear old refresh tokens")
	err = sqlc.Queries.DeleteAuthTokensByUser(dbCtx, refreshToken.UserID, sqldb.AppAuthTokenTypeREFRESH)

	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error clearing old refresh tokens in the database: %s", err)),
		)
		dbSpan.End()
		errMessage := gruenthttp.NewErrorResponse([]any{"Could not store new refresh token"})
		marshalledErrMessage, _ := errMessage.MarshalJson()
		return nil, fiber.NewError(fiber.StatusInternalServerError, string(marshalledErrMessage))
	}

	// Store the new refresh token in the database
	dbCtx, dbSpan = otel.Tracer("sqlc").Start(request.Context, "Store new refresh token")
	_, err = sqlc.Queries.CreateAuthToken(
		dbCtx,
		refreshToken.UserID,
		*newRefreshToken,
		sqldb.AppAuthTokenTypeREFRESH,
		time.Now().Add(gruenthttputil.RefreshTokenExpiry).UTC(),
	)

	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error storing new refresh token in the database: %s", err)),
		)
		dbSpan.End()
		errMessage := gruenthttp.NewErrorResponse([]any{"Could not store new refresh token"})
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
