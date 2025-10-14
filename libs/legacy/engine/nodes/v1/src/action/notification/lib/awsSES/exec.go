package awsSES

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/ses"

	"github.com/aws/aws-sdk-go/aws/session"
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	"github.com/iotea-com/iotea/libs/engine/template"
)

func (n *AwsSESActionSubnode) Exec(params node.ExecParams) node.Error {
Loop:
	for {
		select {
		case <-params.Ctx.Done():
			break Loop
		case ioData := <-n.inputChannels[0].Channel:
			invalid, err := n.handle(ioData)
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

func (n *AwsSESActionSubnode) handle(ioData node.IoData) (bool, error) {
	awsConfig := &aws.Config{
		Region: aws.String(n.config.ThingAwsSES.AwsRegion),
		Credentials: credentials.NewStaticCredentials(
			n.config.ThingAwsSES.AwsAccessKeyId,
			n.config.ThingAwsSES.AwsSecretAccessKey,
			"",
		),
	}

	sender := n.config.Sender
	receiver := n.config.Receiver
	subject := n.config.Subject

	dataMap, err := template.ConvertIoDataToMap(ioData)
	if err != nil {
		return true, fmt.Errorf("failed to convert data to map: %v", err) // Validation error
	}

	message := n.config.TemplateMessage
	if n.config.ModelInput != nil {
		// Preprocess the filter template
		n.template.Preprocess()

		// Fill the message template with data from IoData
		m, err := n.template.Fill(dataMap)
		if err != nil {
			return false, fmt.Errorf("failed to fill template: %v", err)
		}

		message = *m
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return false, fmt.Errorf("failed to create AWS session: %v", err)
	}

	// Create an SES service client
	sesSession := ses.New(sess)

	// Construct the email input parameters
	emailInput := &ses.SendEmailInput{
		Source: aws.String(sender),
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(receiver),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Data: aws.String(message),
				},
			},
			Subject: &ses.Content{
				Data: aws.String(subject),
			},
		},
	}

	n.logger.Info().Ctx(ioData.Ctx).
		Str("data", fmt.Sprintf("%v", emailInput)).
		Msgf("Sending email via SES from %s to %s", n.config.Sender, n.config.Receiver)

	// Send the email
	_, err = sesSession.SendEmailWithContext(ioData.Ctx, emailInput)
	if err != nil {
		return false, fmt.Errorf("failed to send email via SES: %v", err) // Fatal error
	}

	n.logger.Info().Ctx(ioData.Ctx).Msgf("Successfully sent email from %s to %s", n.config.Sender, n.config.Receiver)

	return false, nil // No error
}
