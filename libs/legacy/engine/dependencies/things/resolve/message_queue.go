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

func resolveKafkaProducerDependency(
	env environment.Env,
	producerThing *db.ThingModel,
	nodeConfig map[string]any,
	key string,
	prismaClient *db.PrismaClient,
	secretsClient secrets.SecretsClient,
	spaceId string,
) error {
	_ = env
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

	// Fetch the cluster thing from the database
	clusterThing, err := prismaClient.Thing.FindUnique(db.Thing.ID.Equals(clusterId)).Exec(context.Background())
	if err != nil {
		return fmt.Errorf("an error ocurred whilst fetching Kafka cluster thing \"%v\" from database: %v", clusterId, err)
	}

	// Now we have to resolve the cluster
	var cluster things.KafkaCluster
	err = json.Unmarshal(clusterThing.Attributes, &cluster)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config: %v", err)
	}

	// Set the client's cluster to the fetched cluster attributes
	producer.Cluster = cluster

	// Set the unpacked client attributes to the right field in the node configuration
	nodeConfig[key] = producer

	return nil
}

func resolveKafkaConsumerDependency(
	env environment.Env,
	producerThing *db.ThingModel,
	nodeConfig map[string]any,
	key string,
	prismaClient *db.PrismaClient,
	secretsClient secrets.SecretsClient,
	spaceId string,
) error {
	_ = env
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

	// Fetch the cluster thing from the database
	clusterThing, err := prismaClient.Thing.FindUnique(db.Thing.ID.Equals(clusterId)).Exec(context.Background())
	if err != nil {
		return fmt.Errorf("an error ocurred whilst fetching Kafka cluster thing \"%v\" from database: %v", clusterId, err)
	}

	// Now we have to resolve the cluster
	var cluster things.KafkaCluster
	err = json.Unmarshal(clusterThing.Attributes, &cluster)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config: %v", err)
	}

	// Set the client's cluster to the fetched cluster attributes
	consumer.Cluster = cluster

	// Set the unpacked client attributes to the right field in the node configuration
	nodeConfig[key] = consumer

	return nil
}

func resolveNatsClientDependency(
	env environment.Env,
	clientThing *db.ThingModel,
	nodeConfig map[string]any,
	key string,
	prismaClient *db.PrismaClient,
	secretsClient secrets.SecretsClient,
	spaceId string,
) error {
	_ = env
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

	// Fetch the cluster thing from the database
	serverThing, err := prismaClient.Thing.FindUnique(db.Thing.ID.Equals(serverId)).Exec(context.Background())
	if err != nil {
		return fmt.Errorf("an error ocurred whilst fetching NATS server thing \"%v\" from database: %v", serverId, err)
	}

	// Now we have to resolve the cluster
	var server things.NatsServer
	err = json.Unmarshal(serverThing.Attributes, &server)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config: %v", err)
	}

	// Set the client's server to the fetched broker attributes
	client.Server = server

	// Set the unpacked client attributes to the right field in the node configuration
	nodeConfig[key] = client

	return nil
}
