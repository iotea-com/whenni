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

func resolveMqttClientDependency(
	env environment.Env,
	clientThing *db.ThingModel,
	nodeConfig map[string]any,
	key string,
	prismaClient *db.PrismaClient,
	secretsClient secrets.SecretsClient,
	spaceId string,
) error {
	_ = env

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

	// Fetch the broker thing from the database
	brokerThing, err := prismaClient.Thing.FindUnique(db.Thing.ID.Equals(brokerId)).Exec(context.Background())
	if err != nil {
		return fmt.Errorf("an error ocurred whilst fetching MQTT broker thing \"%s\" from database: %s", brokerId, err)
	}

	// Now we have to resolve the broker
	var broker things.MqttBroker
	err = json.Unmarshal(brokerThing.Attributes, &broker)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config: %s", err)
	}

	// Set the client's broker to the fetched broker attributes
	client.Broker = broker

	// Fetch the client keypair from the secrets engine if its required
	if broker.Protocol == "mqtts" {
		cert, privateKey, err := secretsClient.RetrieveCert(client.CertificateId)
		if err != nil {
			return fmt.Errorf("failed to retrieve the certificate for this MQTT client: %s", err)
		}
		if cert == nil || privateKey == nil {
			return fmt.Errorf("failed to retrieve the certificate for this MQTT client: missing cert or private key")
		}
		client.Keypair.ClientCertPem = *cert
		client.Keypair.ClientCertKey = *privateKey
	}

	// Extract the secret from the space secrets
	if client.Password != "" {
		// Get all the secrets for the space
		spaceSecrets, err := secretsClient.ReadKeyValue(context.Background(), fmt.Sprintf("spaces/%s", spaceId))
		if err != nil {
			return fmt.Errorf("failed to resolve secrets: %v", err)
		}

		client.Password = secretsClient.GetStringFromMap(spaceSecrets, client.Password)
	}

	// Set the unpacked client attributes to the right field in the node configuration
	nodeConfig[key] = client

	return nil
}
