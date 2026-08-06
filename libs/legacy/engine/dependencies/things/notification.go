package things

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/ongruent/gruent/libs/secrets"
)

const AwsSNSThingCategory ThingCategory = "AWS_SNS_ENDPOINT"
const AwsSESThingCategory ThingCategory = "AWS_SES_ENDPOINT"
const SendgridClientThingCategory ThingCategory = "SENDGRID_CLIENT"

type AwsSNS struct {
	AwsAccessKeyId     string `json:"AwsAccessKeyId" validate:"required"`
	AwsSecretAccessKey string `json:"AwsSecretAccessKey" validate:"required"`
	AwsRegion          string `json:"AwsRegion" validate:"required"`
}

func NewAwsSNSFromAttributes(attributes map[string]any, secretsClient secrets.SecretsClient, spaceId string) (*AwsSNS, error) {
	// Unmarshal
	jsonAttributes, _ := json.Marshal(attributes)
	var awsSNS AwsSNS
	err := json.Unmarshal(jsonAttributes, &awsSNS)
	if err != nil {
		return nil, err
	}

	// Get password from secrets engine
	spaceSecrets, err := secretsClient.ReadKeyValue(context.Background(), fmt.Sprintf("spaces/%s", spaceId))
	if err != nil {
		return nil, fmt.Errorf("failed to resolve secrets: %v", err)
	}

	// Extract the secret from the space secrets
	awsSNS.AwsSecretAccessKey = secretsClient.GetStringFromMap(spaceSecrets, awsSNS.AwsSecretAccessKey)

	// Validate
	err = awsSNS.Validate()
	if err != nil {
		return nil, err
	}

	return &awsSNS, nil
}

func (a *AwsSNS) Category() ThingCategory {
	return AwsSNSThingCategory
}

func (a *AwsSNS) Validate() error {
	v := validator.New()

	if err := v.Struct(a); err != nil {
		return err
	}

	return nil
}

func (a *AwsSNS) MarshalJson() ([]byte, error) {
	jsonAttributes, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	return jsonAttributes, nil
}

type AwsSES struct {
	AwsAccessKeyId     string `json:"AwsAccessKeyId" validate:"required"`
	AwsSecretAccessKey string `json:"AwsSecretAccessKey" validate:"required"`
	AwsRegion          string `json:"AwsRegion" validate:"required"`
}

func NewAwsSESFromAttributes(attributes map[string]any, secretsClient secrets.SecretsClient, spaceId string) (*AwsSES, error) {
	// Unmarshal
	jsonAttributes, _ := json.Marshal(attributes)
	var awsSES AwsSES
	err := json.Unmarshal(jsonAttributes, &awsSES)
	if err != nil {
		return nil, err
	}

	// Get password from secrets engine
	spaceSecrets, err := secretsClient.ReadKeyValue(context.Background(), fmt.Sprintf("spaces/%s", spaceId))
	if err != nil {
		return nil, fmt.Errorf("failed to resolve secrets: %v", err)
	}

	// Extract the secret from the space secrets
	awsSES.AwsSecretAccessKey = secretsClient.GetStringFromMap(spaceSecrets, awsSES.AwsSecretAccessKey)

	// Validate
	err = awsSES.Validate()
	if err != nil {
		return nil, err
	}

	return &awsSES, nil
}

func (a *AwsSES) Category() ThingCategory {
	return AwsSESThingCategory
}

func (a *AwsSES) Validate() error {
	v := validator.New()

	if err := v.Struct(a); err != nil {
		return err
	}

	return nil
}

func (a *AwsSES) MarshalJson() ([]byte, error) {
	jsonAttributes, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	return jsonAttributes, nil
}

type SendgridClient struct {
	ApiKey string `json:"apiKey" validate:"required"`
}

func NewSendgridClientFromAttributes(attributes map[string]any, secretsClient secrets.SecretsClient, spaceId string) (*SendgridClient, error) {
	// Unmarshal
	jsonAttributes, _ := json.Marshal(attributes)
	var sendgridClient SendgridClient
	err := json.Unmarshal(jsonAttributes, &sendgridClient)
	if err != nil {
		return nil, err
	}

	// Get password from secrets engine
	spaceSecrets, err := secretsClient.ReadKeyValue(context.Background(), fmt.Sprintf("spaces/%s", spaceId))
	if err != nil {
		return nil, fmt.Errorf("failed to resolve secrets: %v", err)
	}

	// Extract the secret from the space secrets
	sendgridClient.ApiKey = secretsClient.GetStringFromMap(spaceSecrets, sendgridClient.ApiKey)

	// Validate
	err = sendgridClient.Validate()
	if err != nil {
		return nil, err
	}

	return &sendgridClient, nil
}

func (m *SendgridClient) Category() ThingCategory {
	return SendgridClientThingCategory
}

func (m *SendgridClient) Validate() error {
	v := validator.New()

	if err := v.Struct(m); err != nil {
		return err
	}

	return nil
}

func (m *SendgridClient) MarshalJson() ([]byte, error) {
	jsonAttributes, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return jsonAttributes, nil
}
