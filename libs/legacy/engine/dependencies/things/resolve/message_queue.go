package resolve

import (
	"encoding/json"
	"fmt"

	sqldb "github.com/ongruent/gruent/db/sqlc"
	"github.com/ongruent/gruent/libs/legacy/engine/dependencies/things"
	"github.com/ongruent/gruent/libs/legacy/engine/environment"
	"github.com/ongruent/gruent/libs/secrets"
)

func resolveKafkaProducerDependency(
	env environment.Env,
	producerThing sqldb.AppThing,
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
	var producer things.KafkaProducer
	err := json.Unmarshal(producerThing.Attributes, &producer)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config: %v", err)
	}

	// Verify that the broker was passed as the correct type
	clusterId, ok := producer.Cluster.(string)
	if !ok {
		return fmt.Errorf("invalid cluster configuration in database: expected a string but got %T", producer.Cluster)
	}

	return fmt.Errorf("resolving Kafka cluster thing %q requires database access: pending sqlc migration", clusterId)
}

func resolveKafkaConsumerDependency(
	env environment.Env,
	producerThing sqldb.AppThing,
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

	// Unmarshal the thing so that we can process it
	var consumer things.KafkaConsumer
	err := json.Unmarshal(producerThing.Attributes, &consumer)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config: %v", err)
	}

	// Verify that the broker was passed as the correct type
	clusterId, ok := consumer.Cluster.(string)
	if !ok {
		return fmt.Errorf("invalid cluster configuration in database: expected a string but got %T", consumer.Cluster)
	}

	return fmt.Errorf("resolving Kafka cluster thing %q requires database access: pending sqlc migration", clusterId)
}

func resolveNatsClientDependency(
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
	var client things.NatsClient
	err := json.Unmarshal(clientThing.Attributes, &client)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config: %v", err)
	}

	// Verify that the broker was passed as the correct type
	serverId, ok := client.Server.(string)
	if !ok {
		return fmt.Errorf("invalid cluster configuration in database: expected a string but got %T", client.Server)
	}

	return fmt.Errorf("resolving NATS server thing %q requires database access: pending sqlc migration", serverId)
}
