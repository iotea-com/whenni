package organizationsMembersInvite

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	ioteahttputil "github.com/iotea-com/iotea/libs/http/util"
	"github.com/iotea-com/iotea/services/http-api/config"
	"github.com/iotea-com/iotea/services/http-api/services/smtp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func execute(request *ioteahttp.Request[Input]) (*Output, error) {
	request.Span.AddEvent("execute")
	request.Span.SetAttributes(
		attribute.String("request.Input.OrgId", request.Input.OrgId),
		attribute.String("request.Input.Email", request.Input.Email),
		attribute.String("request.Input.AppUrl", request.Input.AppUrl),
	)

	// Create a new invite token
	token, err := ioteahttputil.GenerateInvitationToken(request.Input.OrgId, request.Actor.Id, request.Input.Email, config.VaultConf.JwtSecret)
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "invite_generation"),
			attribute.String("error.message", fmt.Sprintf("error generating invite: %s", err)),
		)
		return nil, fmt.Errorf("error generating invitation token: %s", err)
	}

	// Send an email with the token
	_, smtpSpan := otel.Tracer("smtp").Start(request.Context, "Send invite email")
	from := "auth@iotea.com"

	to := []string{request.Input.Email}

	const htmlBodyTemplate = `
	<html>
		<body>
			<div>
				<h1>IOTEA Invitation</h1>
				<p>You've been invited to join an IOTEA organization!</p>
				<p>
					<a href="{{ .AppUrl }}/signup?inviteToken={{ .Token }}" style="
						display: inline-block;
						background: #0070f3;
						color: white;
						padding: 12px 24px;
						text-decoration: none;
						border-radius: 5px;
					">
						Click here to sign up
					</a>
				</p>
			</div>
		</body>
	</html>
	`

	const plainBodyTemplate = `IOTEA Invitation

	Copy and paste the link below to sign up and join the organization:
	{{ .AppUrl }}/signup?inviteToken={{ .Token }}

	If you did not expect this invitation, please ignore this email.
	This link will expire in {{ .Expiry }} minutes.`

	bodyData := map[string]string{
		"AppUrl": request.Input.AppUrl,
		"Token":  *token,
		"Expiry": fmt.Sprintf("%.0f", ioteahttputil.InvitationTokenExpiry.Minutes()),
	}

	// Create multipart message
	subject := "IOTEA Invitation"
	message, err := ioteahttputil.BuildEmail(from, to, subject, plainBodyTemplate, htmlBodyTemplate, bodyData)
	if err != nil {
		smtpSpan.SetAttributes(
			attribute.String("error.type", "multipart_message_build"),
			attribute.String("error.message", fmt.Sprintf("error building multipart message: %s", err)),
		)
		smtpSpan.End()
		errResponse := ioteahttp.NewErrorResponse([]any{"Unable to send invite email. Please try again."})
		marshalledErrResponse, _ := errResponse.MarshalJson()
		return nil, fiber.NewError(fiber.StatusInternalServerError, string(marshalledErrResponse))
	}

	err = smtp.SmtpClient.SendMail(from, to, subject, message)
	if err != nil {
		smtpSpan.SetAttributes(
			attribute.String("error.type", "smtp_send"),
			attribute.String("error.message", fmt.Sprintf("error sending magic link email: %s", err)),
		)
		smtpSpan.End()
		errResponse := ioteahttp.NewErrorResponse([]any{"Unable to send invite email. Please try again."})
		marshalledErrResponse, _ := errResponse.MarshalJson()
		return nil, fiber.NewError(fiber.StatusInternalServerError, string(marshalledErrResponse))
	}

	smtpSpan.End()

	return nil, nil
}
