package sendgridActionNode

import (
	"fmt"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	"github.com/iotea-com/iotea/libs/engine/template"
)

func (n *SendgridActionSubnode) Exec(params node.ExecParams) node.Error {
	// Client library works with a registered callback. In this loop we simply
	// await for a context cancellation, or an error
Loop:
	for {
		select {
		case <-params.Ctx.Done():
			break Loop
		case ioData := <-n.inputChannels[0].Channel:
			invalid, err := n.handleNotification(ioData)
			if err != nil {
				if invalid {
					// Safely notify the runtime that data was processed
					n.notifyChannel <- node.Notification{
						Type:    node.NotifyDataInvalid,
						DataCtx: ioData.Ctx,
						Reason:  err.Error(),
					}
				} else {
					return node.Error{
						Type:   node.FatalError,
						Reason: err.Error(),
					}
				}
			} else {
				// Safely notify the runtime that data was processed
				n.notifyChannel <- node.Notification{
					Type:    node.NotifyDataProcessed,
					DataCtx: ioData.Ctx,
				}
			}
		}
	}

	return node.Error{
		Type: node.NoError,
	}
}

func (n *SendgridActionSubnode) handleNotification(ioData node.IoData) (bool, error) {
	sendAs := n.config.Sender
	if n.config.SendAs != nil {
		sendAs = *n.config.SendAs
	}
	from := mail.NewEmail(sendAs, n.config.Sender)
	receiver := mail.NewEmail(n.config.Receiver, n.config.Receiver)
	subject := n.config.Subject

	// Fill the message template if the input schema is provided
	htmlContent := n.config.TemplateMessage
	plainTextContent := n.config.TemplateMessage
	if n.config.ModelInput != nil {
		// Preprocess the filter template
		n.template.Preprocess()

		// Fill the message template with data from IoData
		dataMap, err := template.ConvertIoDataToMap(ioData)
		if err != nil {
			return true, fmt.Errorf("failed to convert data to map: %v", err) // Validation error
		}

		ptc, err := n.template.Fill(dataMap)
		if err != nil {
			return true, fmt.Errorf("failed to fill template: %v", err)
		}

		htmlContent = *ptc
		plainTextContent = *ptc
	}

	// Create a Sendgrid service client
	client := sendgrid.NewSendClient(n.config.SendgridClient.ApiKey)

	// Construct the email input parameters
	message := mail.NewSingleEmail(from, subject, receiver, plainTextContent, htmlContent)

	n.logger.Info().Ctx(ioData.Ctx).
		Str("data", fmt.Sprintf("%v", message)).
		Msgf("Sending email via Sendgrid API from %s to %s", n.config.Sender, n.config.Receiver)

	// Send the email
	response, err := client.Send(message)
	if err != nil {
		return false, fmt.Errorf("failed to send email via Sendgrid: %v", err)
	}

	if response.StatusCode == 202 {
		n.logger.Info().Ctx(ioData.Ctx).Msgf("Successfully sent email from %s to %s", from.Address, receiver.Address)
	} else if response.Body != "" {
		n.logger.Warn().Ctx(ioData.Ctx).Msgf("Received %d response from Sendgrid API: %s", response.StatusCode, response.Body)
	} else if response.Body == "" {
		n.logger.Warn().Ctx(ioData.Ctx).Msgf("Received %d response from Sendgrid API", response.StatusCode)
	}

	return false, nil // No error
}
