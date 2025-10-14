package resolve

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/iotea-com/iotea/libs/engine/dependencies/things"
	"github.com/iotea-com/iotea/libs/engine/environment"
	"github.com/iotea-com/iotea/libs/secrets"
	"github.com/iotea-com/iotea/prisma/db"
)

func resolveMinioBucketDependency(
	env environment.Env,
	clientThing *db.ThingModel,
	nodeConfig map[string]any,
	key string,
	prismaClient *db.PrismaClient,
	secretsClient secrets.SecretsClient,
	spaceId string,
) error {
	_ = env
	_ = prismaClient

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
	clientThing *db.ThingModel,
	nodeConfig map[string]any,
	key string,
	prismaClient *db.PrismaClient,
	secretsClient secrets.SecretsClient,
	spaceId string,
) error {
	_ = env
	_ = prismaClient

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
