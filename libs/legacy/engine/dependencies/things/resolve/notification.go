package resolve

import (
	"context"
	"encoding/json"
	"fmt"

	sqldb "github.com/ongruent/gruent/db/sqlc"
	"github.com/ongruent/gruent/libs/legacy/engine/dependencies/things"
	"github.com/ongruent/gruent/libs/legacy/engine/environment"
	"github.com/ongruent/gruent/libs/secrets"
)

func resolveAwsSNSDependency(
	env environment.Env,
	snsThing sqldb.AppThing,
	nodeConfig map[string]any,
	key string,
	secretsClient secrets.SecretsClient,
	spaceId string,
) error {
	_ = env

	// Unmarshal the thing so that we can process it
	var awsSNS things.AwsSNS
	err := json.Unmarshal(snsThing.Attributes, &awsSNS)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config: %v", err)
	}

	// Get all the secrets for the space
	spaceSecrets, err := secretsClient.ReadKeyValue(context.Background(), fmt.Sprintf("spaces/%s", spaceId))
	if err != nil {
		return fmt.Errorf("failed to resolve secrets: %v", err)
	}

	// Extract the secret from the space secrets
	awsSNS.AwsSecretAccessKey = secretsClient.GetStringFromMap(spaceSecrets, awsSNS.AwsSecretAccessKey)

	// Set the unpacked SNS attributes to the right field in the node configuration
	nodeConfig[key] = awsSNS

	return nil
}

func resolveAwsSESDependency(
	env environment.Env,
	sesThing sqldb.AppThing,
	nodeConfig map[string]any,
	key string,
	secretsClient secrets.SecretsClient,
	spaceId string,
) error {
	_ = env

	// Unmarshal the thing so that we can process it
	var awsSES things.AwsSES
	err := json.Unmarshal(sesThing.Attributes, &awsSES)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config: %v", err)
	}

	// Get all the secrets for the space
	spaceSecrets, err := secretsClient.ReadKeyValue(context.Background(), fmt.Sprintf("spaces/%s", spaceId))
	if err != nil {
		return fmt.Errorf("failed to resolve secrets: %v", err)
	}

	// Extract the secret from the space secrets
	awsSES.AwsSecretAccessKey = secretsClient.GetStringFromMap(spaceSecrets, awsSES.AwsSecretAccessKey)

	// Set the unpacked SES attributes to the right field in the node configuration
	nodeConfig[key] = awsSES

	return nil
}

func resolveSendgridClientDependency(
	env environment.Env,
	sendgridClientThing sqldb.AppThing,
	nodeConfig map[string]any,
	key string,
	secretsClient secrets.SecretsClient,
	spaceId string,
) error {
	_ = env

	// Unmarshal the thing so that we can process it
	var sendgridClient things.SendgridClient
	err := json.Unmarshal(sendgridClientThing.Attributes, &sendgridClient)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config: %v", err)
	}

	// Get all the secrets for the space
	spaceSecrets, err := secretsClient.ReadKeyValue(context.Background(), fmt.Sprintf("spaces/%s", spaceId))
	if err != nil {
		return fmt.Errorf("failed to resolve secrets: %v", err)
	}

	// Extract the secret from the space secrets
	sendgridClient.ApiKey = secretsClient.GetStringFromMap(spaceSecrets, sendgridClient.ApiKey)

	// Set the unpacked Sendgrid client attributes to the right field in the node configuration
	nodeConfig[key] = sendgridClient

	return nil
}
