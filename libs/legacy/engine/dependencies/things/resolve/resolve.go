package resolve

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ongruent/gruent/libs/legacy/engine/channels"
	"github.com/ongruent/gruent/libs/legacy/engine/environment"
	"github.com/ongruent/gruent/libs/secrets"
)

const (
	ThingFieldNameSuffix = "::thing"
)

type ResolveConfig struct {
	Env           environment.Env
	SpaceId       string
	Metadata      *channels.NodeMetadata // This is the node metadate we're going to attempt to populate things from the database into
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
			found, err := findAndResolveDependency(config.Env, unmarshalledConfig, &nodeDependency, config.SecretsClient, config.SpaceId)
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
	secretsClient secrets.SecretsClient,
	spaceId string,
) (bool, error) {
	for key, value := range nodeConfig {
		if key == thingDependecy.FieldName {
			if val, ok := value.(string); ok && val == thingDependecy.ThingId {
				if strings.HasSuffix(thingDependecy.FieldName, ThingFieldNameSuffix) {
					return false, fmt.Errorf("thing dependency resolution is not available: pending sqlc migration")
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
