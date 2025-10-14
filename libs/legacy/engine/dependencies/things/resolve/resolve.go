package resolve

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/iotea-com/iotea/libs/engine/channels"
	"github.com/iotea-com/iotea/libs/engine/dependencies/things"
	"github.com/iotea-com/iotea/libs/engine/environment"
	"github.com/iotea-com/iotea/libs/secrets"
	"github.com/iotea-com/iotea/prisma/db"
)

const (
	ThingFieldNameSuffix = "::thing"
)

type ResolveConfig struct {
	Env           environment.Env
	SpaceId       string
	Metadata      *channels.NodeMetadata // This is the node metadate we're going to attempt to populate things from the database into
	PrismaClient  *db.PrismaClient       // We use the database client to fetch the thing from the database table
	SecretsClient secrets.SecretsClient  // Used to fetch certificates and other sensitive information
}

func ResolveDependenciesInNode(config *ResolveConfig) error {
	if len(config.Metadata.Dependencies.Things) > 0 {

		// First we attemp to unmarshall the node configuration into a type we can manipulate
		var unmarshalledConfig map[string]any
		err := json.Unmarshal(config.Metadata.Config, &unmarshalledConfig)
		if err != nil {
			return fmt.Errorf("failed to unmarshal config: %w", err)
		}

		// For each dependency in the node, call the right thing resolve function
		for _, nodeDependency := range config.Metadata.Dependencies.Things {
			found, err := findAndResolveDependency(config.Env, unmarshalledConfig, &nodeDependency, config.PrismaClient, config.SecretsClient, config.SpaceId)
			if err != nil {
				return err
			}
			if !found {
				return fmt.Errorf("dependency check failed for fieldName '%s' with thingId '%s'. Field was in the dependencies object, but was not found in the node configuration object", nodeDependency.FieldName, nodeDependency.ThingId)
			}
		}

		// Re-marshal updated config back to JSON
		updatedConfig, err := json.Marshal(unmarshalledConfig)
		if err != nil {
			return fmt.Errorf("failed to re-marshal updated config: %w", err)
		}
		config.Metadata.Config = updatedConfig
	}

	return nil
}

func findAndResolveDependency(
	env environment.Env,
	nodeConfig map[string]any,
	thingDependecy *channels.NodeThingDependency,
	prismaClient *db.PrismaClient,
	secretsClient secrets.SecretsClient,
	spaceId string,
) (bool, error) {
	for key, value := range nodeConfig {
		if key == thingDependecy.FieldName {
			if val, ok := value.(string); ok && val == thingDependecy.ThingId {
				if strings.HasSuffix(thingDependecy.FieldName, ThingFieldNameSuffix) {
					// Fetch the thing from the database
					thing, err := prismaClient.Thing.FindUnique(db.Thing.ID.Equals(val)).Exec(context.Background())
					if err != nil {
						return false, fmt.Errorf("an error ocurred whilst fetching thing \"%v\" from database: %v", val, err)
					}

					// Resolve based on the category type
					switch thing.ThingCategory {
					case things.HttpServerThingCategory.String():
						if err := resolveHttpServerDependency(env, thing, nodeConfig, key, prismaClient, secretsClient, spaceId); err != nil {
							return false, fmt.Errorf("could not resolve HTTP_SERVER dependency: %v", err)
						}
					case things.MqttClientThingCategory.String():
						if err := resolveMqttClientDependency(env, thing, nodeConfig, key, prismaClient, secretsClient, spaceId); err != nil {
							return false, fmt.Errorf("could not resolve MQTT_CLIENT dependency: %v", err)
						}
					case things.S3BucketThingCategory.String():
						if err := resolveS3BucketDependency(env, thing, nodeConfig, key, prismaClient, secretsClient, spaceId); err != nil {
							return false, fmt.Errorf("could not resolve S3_BUCKET dependency: %v", err)
						}
					case things.MinioBucketThingCategory.String():
						if err := resolveMinioBucketDependency(env, thing, nodeConfig, key, prismaClient, secretsClient, spaceId); err != nil {
							return false, fmt.Errorf("could not resolve MINIO_BUCKET dependency: %v", err)
						}
					case things.KafkaProducerThingCategory.String():
						if err := resolveKafkaProducerDependency(env, thing, nodeConfig, key, prismaClient, secretsClient, spaceId); err != nil {
							return false, fmt.Errorf("could not resolve KAFKA_PRODUCER dependency: %v", err)
						}
					case things.KafkaConsumerThingCategory.String():
						if err := resolveKafkaConsumerDependency(env, thing, nodeConfig, key, prismaClient, secretsClient, spaceId); err != nil {
							return false, fmt.Errorf("could not resolve KAFKA_CONSUMER dependency: %v", err)
						}
					case things.NatsClientThingCategory.String():
						if err := resolveNatsClientDependency(env, thing, nodeConfig, key, prismaClient, secretsClient, spaceId); err != nil {
							return false, fmt.Errorf("could not resolve NATS_CLIENT dependency: %v", err)
						}
					case things.InfluxDbDatabaseThingCategory.String():
						if err := resolveInfluxDbDependency(env, thing, nodeConfig, key, prismaClient, secretsClient, spaceId); err != nil {
							return false, fmt.Errorf("could not resolve INFLUXDB_DATABASE dependency: %v", err)
						}
					case things.ClickhouseDatabaseThingCategory.String():
						if err := resolveClickhouseDatabaseDependency(env, thing, nodeConfig, key, prismaClient, secretsClient, spaceId); err != nil {
							return false, fmt.Errorf("could not resolve CLICKHOUSE_DATABASE dependency: %v", err)
						}
					case things.AwsSNSThingCategory.String():
						if err := resolveAwsSNSDependency(env, thing, nodeConfig, key, prismaClient, secretsClient, spaceId); err != nil {
							return false, fmt.Errorf("could not resolve AWS_SNS_ENDPOINT dependency: %v", err)
						}
					case things.AwsSESThingCategory.String():
						if err := resolveAwsSESDependency(env, thing, nodeConfig, key, prismaClient, secretsClient, spaceId); err != nil {
							return false, fmt.Errorf("could not resolve AWS_SES_ENDPOINT dependency: %v", err)
						}
					case things.SendgridClientThingCategory.String():
						if err := resolveSendgridClientDependency(env, thing, nodeConfig, key, prismaClient, secretsClient, spaceId); err != nil {
							return false, fmt.Errorf("could not resolve SENDGRID_CLIENT dependency: %v", err)
						}
					case things.MongoDbServerThingCategory.String():
						if err := resolveMongoDbDependency(env, thing, nodeConfig, key, prismaClient, secretsClient, spaceId); err != nil {
							return false, fmt.Errorf("could not resolve MONGODB_SERVER dependency: %v", err)
						}
					default:
						return false, fmt.Errorf("unknown thingCategory type %v", thing.ThingCategory)
					}

					// Update the thing in the database after resolving the dependency
					_, err = prismaClient.Thing.FindUnique(db.Thing.ID.Equals(thing.ID)).
						Update(db.Thing.Attributes.Set(thing.Attributes)).
						Exec(context.Background())
					if err != nil {
						return false, fmt.Errorf("failed to update thing in database: %v", err)
					}
					return true, nil // Dependency resolved and thing updated
				} else {
					// A field in the node configuration is requesting a thing, but it does not end the agreed upon suffix
					// for the configuration field name
					return true, fmt.Errorf("node has field %v requesting a thing dependency, but field does NOT have suffix ::thing", thingDependecy.FieldName)
				}
			}
		}
	}

	return false, nil // No matching dependency found
}
