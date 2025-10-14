package resolve

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/iotea-com/iotea/libs/engine/channels"
	"github.com/iotea-com/iotea/prisma/db"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type ExpandConfig struct {
	Metadata     *channels.NodeMetadata // This is the node metadate we're going to attempt to populate things from the database into
	PrismaClient *db.PrismaClient       // We use the database client to fetch the thing from the database table
}

func ExpandInNode(config *ExpandConfig) error {
	if len(config.Metadata.Dependencies.Things) > 0 {
		// Attempt to unmarshal the node configuration into a type we can manipulate
		var unmarshalledConfig map[string]any
		err := json.Unmarshal(config.Metadata.Config, &unmarshalledConfig)
		if err != nil {
			return fmt.Errorf("failed to unmarshal config: %w", err)
		}

		// For each thing entry in the dependencies, verify that it's valid, fetch it from the database
		// and then replace it into the metadata.Config field
		for _, nodeThing := range config.Metadata.Dependencies.Things {
			found, err := findAndUpdateJSON(context.Background(), unmarshalledConfig, nodeThing.FieldName, nodeThing.ThingId, config.PrismaClient)
			if err != nil {
				return err
			}
			if !found {
				return fmt.Errorf("dependency check failed for fieldName '%s' with thingId '%s'", nodeThing.FieldName, nodeThing.ThingId)
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

func findAndUpdateJSON(ctx context.Context, thingDependencies map[string]any, key, value string, prismaClient *db.PrismaClient) (bool, error) {
	for k, v := range thingDependencies {
		if k == key {
			if val, ok := v.(string); ok && val == value {
				return updateFromDatabase(ctx, thingDependencies, key, value, prismaClient)
			}
		}
	}
	return false, nil
}

func updateFromDatabase(ctx context.Context, data map[string]any, key, value string, prismaClient *db.PrismaClient) (bool, error) {
	tracer := otel.Tracer("prisma")
	dbCtx, dbSpan := tracer.Start(ctx, "Get thing")
	defer dbSpan.End()

	thing, err := prismaClient.Thing.FindUnique(
		db.Thing.ID.Equals(value),
	).Exec(dbCtx)
	if err != nil {
		dbSpan.SetAttributes(attribute.String("error", err.Error()))
		return false, err
	}

	var attributes map[string]any
	err = json.Unmarshal([]byte(thing.Attributes), &attributes)
	if err != nil {
		return false, fmt.Errorf("failed to parse attributes: %w", err)
	}

	dbSpan.AddEvent("successfully got thing from the database")

	// Replace the existing value with the fetched, unmarshalled thing
	dbSpan.SetAttributes(attribute.String(fmt.Sprintf("update_from_database_thing_%s.attributes", key), fmt.Sprintf("%v", attributes)))
	data[key] = attributes

	return true, nil
}
