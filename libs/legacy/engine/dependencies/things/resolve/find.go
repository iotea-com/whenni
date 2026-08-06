package resolve

import (
	"context"
	"encoding/json"
	"fmt"

	sqldb "github.com/ongruent/gruent/db/sqlc"
	"github.com/ongruent/gruent/libs/legacy/engine/channels"
)

type FindConfig struct {
	Metadata    *channels.NodeMetadata // This is the node metadate we're going to attempt to populate models from the database into
	SqlcQueries *sqldb.Queries         // We use the database client to fetch the model from the database table
}

func FindDependenciesInNode(config *FindConfig) ([]sqldb.AppThing, error) {
	if len(config.Metadata.Dependencies.Things) > 0 {
		// First we attempt to unmarshal the node configuration into a type we can manipulate
		var unmarshalledConfig map[string]any
		err := json.Unmarshal(config.Metadata.Config, &unmarshalledConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal config: %w", err)
		}

		// For each dependency in the node, find the dependency in the database
		dependencies := []sqldb.AppThing{}
		for _, nodeDependency := range config.Metadata.Dependencies.Things {
			thing, err := findDependency(context.Background(), unmarshalledConfig, &nodeDependency, config.SqlcQueries)
			if err != nil {
				return nil, err
			}

			dependencies = append(dependencies, thing)
		}

		return dependencies, nil
	}

	return []sqldb.AppThing{}, nil
}

func findDependency(ctx context.Context, nodeConfig map[string]any, thingDependecy *channels.NodeThingDependency, sqlcQueries *sqldb.Queries) (sqldb.AppThing, error) {
	for key, value := range nodeConfig {
		if key == thingDependecy.FieldName {
			if val, ok := value.(string); ok && val == thingDependecy.ThingId {
				thing, err := sqlcQueries.GetThing(ctx, val)
				if err != nil {
					return sqldb.AppThing{}, err
				}

				return thing, nil
			}
		}
	}
	return sqldb.AppThing{}, fmt.Errorf("no valid thing found with ID %s", thingDependecy.ThingId)
}
