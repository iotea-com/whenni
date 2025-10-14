package things

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/iotea-com/iotea/libs/secrets"
)

const S3BucketThingCategory ThingCategory = "S3_BUCKET"
const MinioBucketThingCategory ThingCategory = "MINIO_BUCKET"

type S3Bucket struct {
	AwsAccessKeyId     string `json:"AwsAccessKeyId" validate:"required"`
	AwsSecretAccessKey string `json:"AwsSecretAccessKey" validate:"required"`
	AwsRegion          string `json:"AwsRegion" validate:"required"`
	BucketName         string `json:"BucketName" validate:"required"`
}

func NewS3BucketFromAttributes(attributes map[string]any, secretsClient secrets.SecretsClient, spaceId string) (*S3Bucket, error) {
	// Unmarshal
	jsonAttributes, _ := json.Marshal(attributes)
	var s3Bucket S3Bucket
	err := json.Unmarshal(jsonAttributes, &s3Bucket)
	if err != nil {
		return nil, err
	}

	// Get AWS key from secrets engine
	spaceSecrets, err := secretsClient.ReadKeyValue(context.Background(), fmt.Sprintf("spaces/%s", spaceId))
	if err != nil {
		return nil, fmt.Errorf("failed to resolve secrets: %v", err)
	}

	// Extract the secret from the space secrets
	s3Bucket.AwsSecretAccessKey = secretsClient.GetStringFromMap(spaceSecrets, s3Bucket.AwsSecretAccessKey)

	// Validate
	err = s3Bucket.Validate()
	if err != nil {
		return nil, err
	}

	return &s3Bucket, nil
}

func (s *S3Bucket) Category() ThingCategory {
	return S3BucketThingCategory
}

func (s *S3Bucket) Validate() error {
	v := validator.New()

	if err := v.Struct(s); err != nil {
		return err
	}

	return nil
}

func (s *S3Bucket) MarshalJson() ([]byte, error) {
	jsonAttributes, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return jsonAttributes, nil
}

type MinioBucket struct {
	MinioEndpoint        string `json:"MinioEndpoint" validate:"required"`        // The MinIO server URL or IP address
	MinioAccessKeyId     string `json:"MinioAccessKeyId" validate:"required"`     // MinIO Access Key
	MinioSecretAccessKey string `json:"MinioSecretAccessKey" validate:"required"` // MinIO Secret Key
	BucketName           string `json:"BucketName" validate:"required"`           // The name of the MinIO bucket
	UseSSL               bool   `json:"UseSSL"`                                   // Whether to use SSL for the connection (HTTPS)
}

func NewMinioBucketFromAttributes(attributes map[string]any, secretsClient secrets.SecretsClient, spaceId string) (*MinioBucket, error) {
	// Unmarshal
	jsonAttributes, _ := json.Marshal(attributes)
	var minioBucket MinioBucket
	err := json.Unmarshal(jsonAttributes, &minioBucket)
	if err != nil {
		return nil, err
	}

	// Get MinIO key from secrets engine
	spaceSecrets, err := secretsClient.ReadKeyValue(context.Background(), fmt.Sprintf("spaces/%s", spaceId))
	if err != nil {
		return nil, fmt.Errorf("failed to resolve secrets: %v", err)
	}

	// Extract the secret from the space secrets
	minioBucket.MinioSecretAccessKey = secretsClient.GetStringFromMap(spaceSecrets, minioBucket.MinioSecretAccessKey)

	// Validate
	err = minioBucket.Validate()
	if err != nil {
		return nil, err
	}

	return &minioBucket, nil
}

func (m *MinioBucket) Category() ThingCategory {
	return MinioBucketThingCategory
}

func (m *MinioBucket) Validate() error {
	v := validator.New()

	if err := v.Struct(m); err != nil {
		return err
	}

	return nil
}

func (m *MinioBucket) MarshalJson() ([]byte, error) {
	jsonAttributes, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return jsonAttributes, nil
}
