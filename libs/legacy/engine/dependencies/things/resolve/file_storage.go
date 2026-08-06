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

func resolveMinioBucketDependency(
	env environment.Env,
	clientThing sqldb.AppThing,
	nodeConfig map[string]any,
	key string,
	secretsClient secrets.SecretsClient,
	spaceId string,
) error {
	_ = env

	// Unmarshal the thing so that we can process it
	var bucket things.MinioBucket
	err := json.Unmarshal(clientThing.Attributes, &bucket)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config: %v", err)
	}

	// Get all the secrets for the space
	spaceSecrets, err := secretsClient.ReadKeyValue(context.Background(), fmt.Sprintf("spaces/%s", spaceId))
	if err != nil {
		return fmt.Errorf("failed to resolve secrets: %v", err)
	}

	// Extract the secret from the space secrets
	bucket.MinioSecretAccessKey = secretsClient.GetStringFromMap(spaceSecrets, bucket.MinioSecretAccessKey)

	// Set the unpacked bucket attributes to the right field in the node configuration
	nodeConfig[key] = bucket

	return nil
}

func resolveS3BucketDependency(
	env environment.Env,
	clientThing sqldb.AppThing,
	nodeConfig map[string]any,
	key string,
	secretsClient secrets.SecretsClient,
	spaceId string,
) error {
	_ = env

	// Unmarshal the thing so that we can process it
	var bucket things.S3Bucket
	err := json.Unmarshal(clientThing.Attributes, &bucket)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config: %v", err)
	}

	// Get all the secrets for the space
	spaceSecrets, err := secretsClient.ReadKeyValue(context.Background(), fmt.Sprintf("spaces/%s", spaceId))
	if err != nil {
		return fmt.Errorf("failed to resolve secrets: %v", err)
	}

	// Extract the secret from the space secrets
	bucket.AwsSecretAccessKey = secretsClient.GetStringFromMap(spaceSecrets, bucket.AwsSecretAccessKey)

	// Set the unpacked bucket attributes to the right field in the node configuration
	nodeConfig[key] = bucket

	return nil
}
