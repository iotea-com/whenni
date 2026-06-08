package resolve

import (
	"encoding/json"
	"fmt"

	sqldb "github.com/iotea-com/iotea/db/sqlc"
	"github.com/iotea-com/iotea/libs/legacy/engine/dependencies/things"
	"github.com/iotea-com/iotea/libs/legacy/engine/environment"
	"github.com/iotea-com/iotea/libs/secrets"
)

func resolveMqttClientDependency(
	env environment.Env,
	clientThing sqldb.AppThing,
	nodeConfig map[string]any,
	key string,
	secretsClient secrets.SecretsClient,
	spaceId string,
) error {
	_ = env
	_ = nodeConfig
	_ = key
	_ = secretsClient
	_ = spaceId

	// Unmarshall the thing so that we can process it
	var client things.MqttClient
	err := json.Unmarshal(clientThing.Attributes, &client)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config: %s", err)
	}

	// Verify that the broker was passed as the correct type
	brokerId, ok := client.Broker.(string)
	if !ok {
		return fmt.Errorf("invalid broker configuration in database: expected a string but got %T", client.Broker)
	}

	return fmt.Errorf("resolving MQTT broker thing %q requires database access: pending sqlc migration", brokerId)
}
