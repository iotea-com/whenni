package signup

import (
	"context"
	// "fmt"
	// "net/url"
	// "strings"
	// "time"
	// ioteahttputil "github.com/iotea-com/iotea/libs/http/util"
	// "github.com/iotea-com/iotea/prisma/db"
	// "github.com/iotea-com/iotea/services/http-api/services/prisma"
	// "github.com/iotea-com/iotea/services/http-api/services/smtp"
	// "github.com/iotea-com/iotea/services/http-api/util"
	// "go.opentelemetry.io/otel"
	// "go.opentelemetry.io/otel/attribute"
)

// const (
// 	VerificationTokenWordCount     = 4
// 	VerificationTokenWordMinLength = 4
// 	VerificationTokenWordMaxLength = 4
// 	VerificationTokenExpiry        = time.Minute * 15
// )

func action(input *Input, ctx context.Context) error {
	// return sendEmailVerificationEmail(input.Email, input.AppUrl, ctx)
	return nil
}

// func sendEmailVerificationEmail(email string, appUrl string, ctx context.Context) error {
// 	// Clear all existing verification tokens for this email
// 	_, dbSpan := otel.Tracer("prisma").Start(ctx, "Delete existing verification tokens")
// 	_, err := prisma.Client.AuthToken.FindMany(
// 		db.AuthToken.User.Where(db.User.Email.Equals(email)),
// 		db.AuthToken.Type.Equals(db.AuthTokenTypeEmailVerification),
// 	).Delete().Exec(ctx)
// 	if err != nil {
// 		dbSpan.SetAttributes(
// 			attribute.String("action.error.type", "database"),
// 			attribute.String("action.error.message", fmt.Sprintf("error deleting existing verification tokens: %s", err)),
// 		)
// 		return err
// 	}
// 	dbSpan.End()

// 	// Generate words for token
// 	uniqueWords := util.GetUniqueRandomWords(util.DefaultWordList, VerificationTokenWordCount, VerificationTokenWordMinLength, VerificationTokenWordMaxLength)

// 	// Join words into a single string
// 	token := strings.Join(uniqueWords, "-")

// 	// Create a new expiration date
// 	expiresAt := time.Now().Add(VerificationTokenExpiry)

// 	// Store new token in database
// 	_, dbSpan = otel.Tracer("prisma").Start(ctx, "Create verification token")
// 	_, err = prisma.Client.AuthToken.CreateOne(
// 		db.AuthToken.Token.Set(token),
// 		db.AuthToken.Type.Set(db.AuthTokenTypeEmailVerification),
// 		db.AuthToken.ExpiresAt.Set(expiresAt.UTC()),
// 		db.AuthToken.User.Link(db.User.Email.Equals(email)),
// 	).Exec(ctx)
// 	if err != nil {
// 		dbSpan.SetAttributes(
// 			attribute.String("action.error.type", "database"),
// 			attribute.String("action.error.message", fmt.Sprintf("error creating verification token: %s", err)),
// 		)
// 		return err
// 	}

// 	dbSpan.End()

// 	// Send an email with the token
// 	_, smtpSpan := otel.Tracer("smtp").Start(ctx, "Send verification email")
// 	from := "auth@iotea.com"

// 	to := []string{email}

// 	const htmlBodyTemplate = `
// 	<html>
// 		<body>
// 			<h1>Verify your email address</h1>
// 			<p><bold>Your token is: {{ .Token }}</bold></p>
// 			<p>Or click the button below to verify your email address:</p>
// 			<a href="{{ .ConfirmationLink }}" style="
// 				display: inline-block;
// 				background: #0070f3;
// 				color: white;
// 				padding: 12px 24px;
// 				text-decoration: none;
// 				border-radius: 5px;
// 			">
// 				Verify Email
// 			</a>
// 			<p>If you did not create this account, please ignore this email.</p>
// 			<p>This link will expire in {{ .Expiry }} minutes.</p>
// 		</body>
// 	</html>
// 	`

// 	const plainBodyTemplate = `Verify your email address

// 	Copy and paste the link below to verify your email address:
// 	{{ .ConfirmationLink }}

// 	If you did not request this verification, please ignore this email.
// 	This link will expire in {{ .Expiry }} minutes.`

// 	confirmationLink := fmt.Sprintf(`%s/api/auth/verify-email?email=%s&token=%s`, appUrl, url.QueryEscape(email), url.QueryEscape(token))
// 	bodyData := map[string]string{
// 		"Token":            token,
// 		"Expiry":           fmt.Sprintf("%.0f", VerificationTokenExpiry.Minutes()),
// 		"ConfirmationLink": confirmationLink,
// 	}

// 	// Create multipart message
// 	subject := "Verify your email address"
// 	message, err := ioteahttputil.BuildEmail(from, to, subject, plainBodyTemplate, htmlBodyTemplate, bodyData)
// 	if err != nil {
// 		smtpSpan.SetAttributes(
// 			attribute.String("action.error.type", "multipart_message_build"),
// 			attribute.String("action.error.message", fmt.Sprintf("error building multipart message: %s", err)),
// 		)
// 		smtpSpan.End()
// 		return err
// 	}

// 	err = smtp.SmtpClient.SendMail(from, to, subject, message)
// 	if err != nil {
// 		smtpSpan.SetAttributes(
// 			attribute.String("action.error.type", "smtp_send"),
// 			attribute.String("action.error.message", fmt.Sprintf("error sending magic link email: %s", err)),
// 		)
// 		smtpSpan.End()
// 		return err
// 	}

// 	smtpSpan.End()
// 	return nil
// }
