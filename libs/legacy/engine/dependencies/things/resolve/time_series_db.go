package resolve

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/iotea-com/iotea/libs/legacy/engine/dependencies/things"
	"github.com/iotea-com/iotea/libs/legacy/engine/environment"
	"github.com/iotea-com/iotea/libs/secrets"
	"github.com/iotea-com/iotea/prisma/db"
)

func resolveInfluxDbDependency(
	env environment.Env,
	databaseThing *db.ThingModel,
	nodeConfig map[string]any,
	key string,
	prismaClient *db.PrismaClient,
	secretsClient secrets.SecretsClient,
	spaceId string,
) error {
	_ = env
	_ = prismaClient

	// Unmarshall the thing so that we can process it
	var database things.InfluxDbDatabase
	err := json.Unmarshal(databaseThing.Attributes, &database)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config: %v", err)
	}

	// Get all the secrets for the space
	spaceSecrets, err := secretsClient.ReadKeyValue(context.Background(), fmt.Sprintf("spaces/%s", spaceId))
	if err != nil {
		return fmt.Errorf("failed to resolve secrets: %v", err)
	}

	// Extract the secret from the space secrets
	database.Token = secretsClient.GetStringFromMap(spaceSecrets, database.Token)

	// Set the unpacked client attributes to the right field in the node configuration
	nodeConfig[key] = database

	return nil
}

func resolveClickhouseDatabaseDependency(
	env environment.Env,
	databaseThing *db.ThingModel,
	nodeConfig map[string]any,
	key string,
	prismaClient *db.PrismaClient,
	secretsClient secrets.SecretsClient,
	spaceId string,
) error {
	_ = env
	_ = prismaClient

	// Unmarshall the thing so that we can process it
	var database things.ClickhouseDatabase
	err := json.Unmarshal(databaseThing.Attributes, &database)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config: %v", err)
	}

	// Get all the secrets for the space
	spaceSecrets, err := secretsClient.ReadKeyValue(context.Background(), fmt.Sprintf("spaces/%s", spaceId))
	if err != nil {
		return fmt.Errorf("failed to resolve secrets: %v", err)
	}

	// Extract the secret from the space secrets
	database.Password = secretsClient.GetStringFromMap(spaceSecrets, database.Password)

	// Set the unpacked client attributes to the right field in the node configuration
	nodeConfig[key] = database

	return nil
}
