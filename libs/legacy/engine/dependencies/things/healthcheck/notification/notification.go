package notificationHealthcheck

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/iotea-com/iotea/libs/engine/dependencies/things"
	"github.com/sendgrid/sendgrid-go"
)

func AwsSNS(attrs *things.AwsSNS) error {
	// Initialize AWS SNS client
	awsConfig := &aws.Config{
		Region: aws.String(attrs.AwsRegion),
		Credentials: credentials.NewStaticCredentials(
			attrs.AwsAccessKeyId,
			attrs.AwsSecretAccessKey,
			"",
		),
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return fmt.Errorf("failed to create AWS session: %v", err)
	}

	snsSession := sns.New(sess)

	// Check if the SNS client is able to retrieve platform application attributes
	_, err = snsSession.GetPlatformApplicationAttributes(&sns.GetPlatformApplicationAttributesInput{})
	if err != nil {
		return fmt.Errorf("AWS SNS health check failed: %v", err)
	}

	return nil
}

func AwsSES(attrs *things.AwsSES) error {
	// Initialize AWS SES client
	awsConfig := &aws.Config{
		Region: aws.String(attrs.AwsRegion),
		Credentials: credentials.NewStaticCredentials(
			attrs.AwsAccessKeyId,
			attrs.AwsSecretAccessKey,
			"",
		),
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return fmt.Errorf("failed to create AWS session: %v", err)
	}

	sesSession := ses.New(sess)

	// Check if the SES client is able to retrieve send quota
	_, err = sesSession.GetSendQuota(&ses.GetSendQuotaInput{})
	if err != nil {
		return fmt.Errorf("AWS SES health check failed: %v", err)
	}

	return nil
}

func SendgridClient(attrs *things.SendgridClient) error {
	// Send a test request to the SendGrid API
	request := sendgrid.GetRequest(attrs.ApiKey, "/v3/scopes", "https://api.sendgrid.com")
	request.Method = "GET"

	_, err := sendgrid.API(request)
	if err != nil {
		return fmt.Errorf("SendGrid health check failed: %v", err)
	}

	return nil
}
