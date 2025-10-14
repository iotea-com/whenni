package awsSNS

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	"github.com/iotea-com/iotea/libs/engine/template"
)

func (n *AwsSNSActionSubnode) Exec(params node.ExecParams) node.Error {
	// Client library works with a registered callback. In this loop we simply
	// await for a context cancellation, or an error
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

func (n *AwsSNSActionSubnode) handle(ioData node.IoData) (bool, error) {
	awsConfig := &aws.Config{
		Region: aws.String(n.config.ThingAwsSNS.AwsRegion),
		Credentials: credentials.NewStaticCredentials(
			n.config.ThingAwsSNS.AwsAccessKeyId,
			n.config.ThingAwsSNS.AwsSecretAccessKey,
			"",
		),
	}

	topicARN := n.config.TopicARN

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
			return true, fmt.Errorf("failed to fill template: %v", err)
		}

		message = *m
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return false, fmt.Errorf("failed to create AWS session: %v", err) // Fatal error
	}

	n.logger.Info().Ctx(ioData.Ctx).
		Str("data", fmt.Sprintf("%v", message)).
		Msgf("Sending message to SNS topic %s", n.config.TopicARN)

	// Create SNS client
	svc := sns.New(sess)

	// Publish the message to the specified SNS topic
	_, err = svc.Publish(&sns.PublishInput{
		Message:  aws.String(message),
		TopicArn: aws.String(topicARN),
	})

	if err != nil {
		return false, fmt.Errorf("failed to publish message to SNS: %v", err) // Fatal error
	}

	n.logger.Info().Ctx(ioData.Ctx).Msgf("Successfully sent message to SNS topic %s", n.config.TopicARN)

	return false, nil // Data processed
}
