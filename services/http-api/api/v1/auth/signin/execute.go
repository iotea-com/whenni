package signin

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	sqldb "github.com/iotea-com/iotea/db/sqlc"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	ioteahttputil "github.com/iotea-com/iotea/libs/http/util"
	"github.com/iotea-com/iotea/services/http-api/config"
	"github.com/iotea-com/iotea/services/http-api/services/smtp"
	"github.com/iotea-com/iotea/services/http-api/services/sqlc"
	"github.com/jackc/pgx/v5"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"golang.org/x/crypto/bcrypt"
)

const (
	MagicLinkTokenLength = 37
	MagicLinkTokenExpiry = time.Minute * 15
)

func execute(request *ioteahttp.Request[Input]) (*Output, error) {
	request.Span.AddEvent("execute")
	request.Span.SetAttributes(
		attribute.String("request.Input.Email", request.Input.Email),
		attribute.String("request.Input.Method", request.Input.Method),
	)

	// Get the user from the database
	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Get user")
	userEmail := request.Input.Email
	user, err := sqlc.Queries.GetUserByEmail(dbCtx, &userEmail)

	if err != nil {
		if err == pgx.ErrNoRows {
			dbSpan.SetAttributes(
				attribute.String("error.type", "user_not_found"),
				attribute.String("error.message", fmt.Sprintf("user with email %s not found", request.Input.Email)),
			)
			dbSpan.End()
			return nil, fiber.NewError(fiber.StatusPreconditionFailed, "No user was found with the provided email")
		}

		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error getting user from the database: %s", err)),
		)
		dbSpan.End()
		return nil, fiber.NewError(fiber.StatusInternalServerError)
	}

	dbSpan.End()

	// Handle login based on method
	switch request.Input.Method {
	case "magicLink":
		err = handleMagicLinkLogin(user, request.Input.AppUrl, request.Context)
		if err != nil {
			request.Span.SetAttributes(
				attribute.String("error.type", "magic_link"),
				attribute.String("error.message", err.Error()),
			)
			errResponse := ioteahttp.NewErrorResponse([]any{err.Error()})
			marshalledErrResponse, _ := errResponse.MarshalJson()
			return nil, fiber.NewError(fiber.StatusInternalServerError, string(marshalledErrResponse))
		}

		return nil, nil
	case "credentials":
		accessToken, refreshToken, err := handleCredentialsLogin(user, request.Input.Password, request.Context)
		if err != nil {
			request.Span.SetAttributes(
				attribute.String("error.type", "credentials"),
				attribute.String("error.message", err.Error()),
			)
			errResponse := ioteahttp.NewErrorResponse([]any{err.Error()})
			marshalledErrResponse, _ := errResponse.MarshalJson()
			return nil, fiber.NewError(fiber.StatusUnauthorized, string(marshalledErrResponse))
		}

		output := &Output{
			AccessToken:  *accessToken,
			RefreshToken: *refreshToken,
		}

		return output, nil
	}

	dbSpan.SetAttributes(
		attribute.String("error.type", "invalid_method"),
		attribute.String("error.message", fmt.Sprintf("invalid login method: %s", request.Input.Method)),
	)
	return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid login method")
}

func handleMagicLinkLogin(user sqldb.AppUser, appUrl string, ctx context.Context) error {
	// Create a new magic link token
	_, tokenSpan := otel.Tracer("prisma").Start(ctx, "Generate magic link token")
	if user.Email == nil || *user.Email == "" {
		tokenSpan.SetAttributes(
			attribute.String("error.type", "user_no_email"),
			attribute.String("error.message", fmt.Sprintf("user '%s' has no email address", user.ID)),
		)
		tokenSpan.End()
		return errors.New("user has no email address")
	}
	token, err := gonanoid.New(MagicLinkTokenLength)
	if err != nil {
		tokenSpan.SetAttributes(
			attribute.String("error.type", "token_generation"),
			attribute.String("error.message", fmt.Sprintf("error generating magic link token: %s", err)),
		)
		tokenSpan.End()
		return errors.New("error generating magic link token")
	}

	// Create a new expiration date
	expiresAt := time.Now().Add(MagicLinkTokenExpiry)
	tokenSpan.End()

	// Store the token in the database
	dbCtx, dbSpan := otel.Tracer("prisma").Start(ctx, "Store magic link token")
	_, err = sqlc.Queries.CreateAuthToken(dbCtx, user.ID, token, sqldb.AppAuthTokenTypeMAGICLINK, expiresAt.UTC())

	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error getting user from the database: %s", err)),
		)
		dbSpan.End()
		return errors.New("error storing magic link token")
	}

	dbSpan.End()

	// Send an email with the token
	_, smtpSpan := otel.Tracer("smtp").Start(ctx, "Send magic link email")
	from := "auth@iotea.com"

	to := []string{*user.Email}

	const htmlBodyTemplate = `
	<html>
		<body>
			<h1>IOTEA Magic Link Sign In</h1>
			<p>Click the link below to login to your account:</p>
			<p>
				<a href="{{ .AppUrl }}/api/magic-link/verify?t={{ .Token }}">Sign in</a>
			</p>
			<p>If you did not request this login, please ignore this email.</p>
			<p>This link will expire in {{ .Expiry }} minutes.</p>
	  </body>
	</html>
	`

	const plainBodyTemplate = `IOTEA Magic Link Sign In

	Copy and paste the link below to login to your account:
	{{ .AppUrl }}/api/magic-link/verify?t={{ .Token }}

	If you did not request this login, please ignore this email.
	This link will expire in {{ .Expiry }} minutes.`

	bodyData := map[string]string{
		"AppUrl": appUrl,
		"Token":  token,
		"Expiry": fmt.Sprintf("%.0f", MagicLinkTokenExpiry.Minutes()),
	}

	// Create multipart message
	subject := "IOTEA Magic Link Sign In"
	message, err := ioteahttputil.BuildEmail(from, to, subject, plainBodyTemplate, htmlBodyTemplate, bodyData)
	if err != nil {
		smtpSpan.SetAttributes(
			attribute.String("error.type", "multipart_message_build"),
			attribute.String("error.message", fmt.Sprintf("error building multipart message: %s", err)),
		)
		smtpSpan.End()
		return errors.New("error sending magic link email")
	}

	err = smtp.SmtpClient.SendMail(from, to, subject, message)
	if err != nil {
		smtpSpan.SetAttributes(
			attribute.String("error.type", "smtp_send"),
			attribute.String("error.message", fmt.Sprintf("error sending magic link email: %s", err)),
		)
		smtpSpan.End()
		return errors.New("error sending magic link email")
	}

	smtpSpan.End()
	return nil
}

func handleCredentialsLogin(user sqldb.AppUser, providedPassword string, ctx context.Context) (*string, *string, error) {
	// Check if email is verified
	if !user.EmailVerified.Valid {
		return nil, nil, errors.New("email is not verified")
	}

	// Match the provided password with the user's password
	if user.Password == nil {
		return nil, nil, errors.New("user password does not have a password")
	}

	err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(providedPassword))
	if err != nil {
		return nil, nil, errors.New("incorrect password")
	}

	// Generate an access token
	_, tokenSpan := otel.Tracer("prisma").Start(ctx, "Generate tokens")
	accessToken, err := ioteahttputil.GenerateAccessToken(user.ID, config.VaultConf.JwtSecret)
	if err != nil {
		tokenSpan.SetAttributes(
			attribute.String("error.type", "jwt_generation"),
			attribute.String("error.message", fmt.Sprintf("error generating jwt: %s", err)),
		)
		tokenSpan.End()
		return nil, nil, errors.New("error generating access token")
	}

	// Generate a refresh token
	refreshToken, err := ioteahttputil.GenerateRefreshToken()
	if err != nil {
		tokenSpan.SetAttributes(
			attribute.String("error.type", "refresh_token_generation"),
			attribute.String("error.message", fmt.Sprintf("error generating refresh token: %s", err)),
		)
		tokenSpan.End()
		return nil, nil, errors.New("error generating refresh token")
	}

	tokenSpan.End()

	// Clear old refresh tokens from the database
	dbCtx, dbSpan := otel.Tracer("prisma").Start(ctx, "Clear old refresh tokens")
	err = sqlc.Queries.DeleteAuthTokensByUser(dbCtx, user.ID, sqldb.AppAuthTokenTypeREFRESH)

	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error clearing old refresh tokens in the database: %s", err)),
		)
		dbSpan.End()
		return nil, nil, errors.New("error storing refresh token")
	}

	dbSpan.End()

	// Store the refresh token in the database
	dbCtx, dbSpan = otel.Tracer("prisma").Start(ctx, "Store refresh token")
	_, err = sqlc.Queries.CreateAuthToken(
		dbCtx,
		user.ID,
		*refreshToken,
		sqldb.AppAuthTokenTypeREFRESH,
		time.Now().Add(ioteahttputil.RefreshTokenExpiry).UTC(),
	)

	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error storing refresh token: %s", err)),
		)
		dbSpan.End()
		return nil, nil, errors.New("error storing refresh token")
	}

	dbSpan.End()

	return accessToken, refreshToken, nil
}
