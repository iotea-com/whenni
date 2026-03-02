package resolve

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/iotea-com/iotea/libs/legacy/engine/channels"
	"github.com/iotea-com/iotea/prisma/db"
)

type FindConfig struct {
	Metadata     *channels.NodeMetadata // This is the node metadate we're going to attempt to populate models from the database into
	PrismaClient *db.PrismaClient       // We use the database client to fetch the model from the database table
}

func FindDependenciesInNode(config *FindConfig) ([]*db.ThingModel, error) {
	if len(config.Metadata.Dependencies.Things) > 0 {
		// First we attempt to unmarshal the node configuration into a type we can manipulate
		var unmarshalledConfig map[string]any
		err := json.Unmarshal(config.Metadata.Config, &unmarshalledConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal config: %w", err)
		}

		// For each dependency in the node, find the dependency in the database
		dependencies := []*db.ThingModel{}
		for _, nodeDependency := range config.Metadata.Dependencies.Things {
			thing, err := findDependency(unmarshalledConfig, &nodeDependency, config.PrismaClient)
			if err != nil {
				return nil, err
			}

			dependencies = append(dependencies, thing)
		}

		return dependencies, nil
	}

	return []*db.ThingModel{}, nil
}

func findDependency(nodeConfig map[string]any, thingDependecy *channels.NodeThingDependency, prismaClient *db.PrismaClient) (*db.ThingModel, error) {
	for key, value := range nodeConfig {
		if key == thingDependecy.FieldName {
			if val, ok := value.(string); ok && val == thingDependecy.ThingId {
				// Fetch the thing from the database
				thing, err := prismaClient.Thing.FindUnique(db.Thing.ID.Equals(val)).Exec(context.Background())
				if err != nil {
					return nil, err
				}

				return thing, nil
			}
		}
	}
	return nil, fmt.Errorf("no valid thing found with ID %s", thingDependecy.ThingId)
}
