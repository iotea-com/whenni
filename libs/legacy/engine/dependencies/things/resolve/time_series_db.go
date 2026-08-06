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

func resolveInfluxDbDependency(
	env environment.Env,
	databaseThing sqldb.AppThing,
	nodeConfig map[string]any,
	key string,
	secretsClient secrets.SecretsClient,
	spaceId string,
) error {
	_ = env

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
	databaseThing sqldb.AppThing,
	nodeConfig map[string]any,
	key string,
	secretsClient secrets.SecretsClient,
	spaceId string,
) error {
	_ = env

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
