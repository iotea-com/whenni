package middleware

import (
	"context"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
	"github.com/golang-jwt/jwt/v5"
	gruenthttputil "github.com/ongruent/gruent/libs/http/util"
	"github.com/ongruent/gruent/services/http-api/config"
	"github.com/ongruent/gruent/services/http-api/services/sqlc"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/trace"
)

var keyauthConfig = keyauth.Config{
	Validator: func(ctx *fiber.Ctx, token string) (bool, error) {
		requestSpan := trace.SpanFromContext(ctx.UserContext())

		// Check if token is a JWT
		t, err := jwt.Parse(token, gruenthttputil.ParseJwt(config.VaultConf.JwtSecret))
		if err != nil {
			requestSpan.AddEvent(fmt.Sprintf("error validating token as JWT in authorization header: %s", err))
		} else {
			// Pass if no errors
			userId, err := t.Claims.GetSubject()
			if err != nil {
				requestSpan.AddEvent(fmt.Sprintf("error getting subject from JWT: %s", err))
				return false, err
			}

			ctx.Locals("user_id", userId)

			requestSpan.AddEvent("successfully validated JWT in authorization header")
			return true, nil
		}

		// Check DB to see if token matches any API keys
		dbCtx := context.Background()
		apiKey, err := sqlc.Queries.GetApiKey(dbCtx, token)

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				requestSpan.AddEvent("API key not found in database")
				return false, keyauth.ErrMissingOrMalformedAPIKey
			}

			requestSpan.AddEvent(fmt.Sprintf("error validating API key: %s", err))
			return false, err
		}

		ctx.Locals("api_key_id", apiKey.ID)
		ctx.Locals("api_key_org_id", apiKey.OrganizationID)

		requestSpan.AddEvent("successfully validated API key")
		return true, nil
	},
	SuccessHandler: func(ctx *fiber.Ctx) error {
		return ctx.Next()
	},
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		requestSpan := trace.SpanFromContext(ctx.UserContext())
		requestSpan.AddEvent(fmt.Sprintf("failed to validate bearer token: %s", err))
		return fiber.NewError(fiber.StatusUnauthorized, "invalid or expired API key")
	},
	KeyLookup:  "header:" + fiber.HeaderAuthorization,
	AuthScheme: "Bearer",
	ContextKey: "token",
}

var ValidateBearerToken = keyauth.New(keyauthConfig)
