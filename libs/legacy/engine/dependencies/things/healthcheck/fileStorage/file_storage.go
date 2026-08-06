package fileStorageHealthcheck

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/ongruent/gruent/libs/legacy/engine/dependencies/things"
	"github.com/minio/minio-go/v7"
	minioCredentials "github.com/minio/minio-go/v7/pkg/credentials"
)

func S3Bucket(attrs *things.S3Bucket) error {
	// Initialize AWS S3 client
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
		return fmt.Errorf("failed to create AWS session: %s", err)
	}

	s3Session := s3.New(sess)

	// Check if the bucket exists
	_, err = s3Session.GetBucketLocation(&s3.GetBucketLocationInput{
		Bucket: aws.String(attrs.BucketName),
	})
	if err != nil {
		return fmt.Errorf("AWS S3 health check failed: %s", err)
	}

	return nil
}

func MinioBucket(attrs *things.MinioBucket) error {
	// Initialize MinIO client
	client, err := minio.New(attrs.MinioEndpoint, &minio.Options{
		Creds:  minioCredentials.NewStaticV4(attrs.MinioAccessKeyId, attrs.MinioSecretAccessKey, ""),
		Secure: attrs.UseSSL,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize MinIO client: %s", err)
	}

	// Check if the bucket exists
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	exists, err := client.BucketExists(ctx, attrs.BucketName)
	if err != nil {
		return fmt.Errorf("MinIO bucket health check failed: %s", err)
	}

	if !exists {
		return fmt.Errorf("MinIO bucket '%s' does not exist", attrs.BucketName)
	}

	return nil
}
