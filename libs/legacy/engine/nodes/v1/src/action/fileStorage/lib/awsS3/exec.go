package awsS3

import (
	"bytes"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	nodeHelpers "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/action/fileStorage/lib/shared"
)

func (n *S3ActionSubnode) Exec(params node.ExecParams) node.Error {
Loop:
	for {
		select {
		case <-params.Ctx.Done():
			break Loop
		case ioData := <-n.inputChannels[0].Channel:
			if err := n.handle(ioData); err != nil {
				return node.Error{
					Type:   node.FatalError,
					Reason: err.Error(),
				}
			} else {
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

func (n *S3ActionSubnode) handle(ioData node.IoData) error {
	// Attempt to detect content type
	data, contentType, err := nodeHelpers.DetectDataContentType(ioData)
	if err != nil {
		return fmt.Errorf("unable to detect data content type: %v", err)
	}

	// Attempt to upload to S3
	if err := n.writeDataToS3(ioData.Ctx, data, contentType); err != nil {
		return fmt.Errorf("unable to upload data to S3: %v", err)
	}

	return nil
}

func (n *S3ActionSubnode) writeDataToS3(dataCtx context.Context, data []byte, contentType string) error {
	awsConfig := &aws.Config{
		Region: aws.String(n.config.ThingAwsS3Bucket.AwsRegion),
		Credentials: credentials.NewStaticCredentials(
			n.config.ThingAwsS3Bucket.AwsAccessKeyId,
			n.config.ThingAwsS3Bucket.AwsSecretAccessKey,
			"",
		),
	}

	// Preprocess the filter template
	n.keyTemplate.Preprocess()

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return fmt.Errorf("failed to create AWS session: %v", err)
	}

	// Create an S3 service client
	svc := s3.New(sess)

	// Prepare the S3 upload input
	input := &s3.PutObjectInput{
		Bucket:      aws.String(n.config.ThingAwsS3Bucket.BucketName),
		Key:         aws.String(n.keyTemplate.Value),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	}

	n.logger.Info().Ctx(dataCtx).Msgf("Uploading to S3 bucket %s with key %s", n.config.ThingAwsS3Bucket.BucketName, n.keyTemplate.Value)

	// Upload the data to S3
	_, err = svc.PutObject(input)
	if err != nil {
		return fmt.Errorf("failed to upload data to S3: %v", err)
	}

	// Log success
	n.logger.Info().Ctx(dataCtx).Msgf("Data successfully uploaded to S3 bucket %s with key %s", n.config.ThingAwsS3Bucket.BucketName, n.keyTemplate.Value)

	return nil
}
