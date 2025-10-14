package resolve

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/iotea-com/iotea/libs/engine/channels"
	"github.com/iotea-com/iotea/libs/engine/dependencies/models"
	"github.com/iotea-com/iotea/prisma/db"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type ResolveConfig struct {
	Metadata     *channels.NodeMetadata // This is the node metadate we're going to attempt to populate models from the database into
	PrismaClient *db.PrismaClient       // We use the database client to fetch the model from the database table
}

func ResolveInNode(config *ResolveConfig) error {
	if len(config.Metadata.Dependencies.Models) > 0 {
		// Attempt to unmarshal the node configuration into a type we can manipulate
		var unmarshalledConfig map[string]any
		err := json.Unmarshal(config.Metadata.Config, &unmarshalledConfig)
		if err != nil {
			return fmt.Errorf("failed to unmarshal config: %w", err)
		}

		// For each model entry in the dependencies, verify that it's valid, fetch it from the database
		// and then replace it into the metadata.Config field
		for _, nodeModel := range config.Metadata.Dependencies.Models {
			found, err := findAndUpdateJSON(context.Background(), unmarshalledConfig, nodeModel.FieldName, nodeModel.ModelId, config.PrismaClient)
			if err != nil {
				return err
			}
			if !found {
				return fmt.Errorf("dependency check failed for fieldName '%s' with modelId '%s'", nodeModel.FieldName, nodeModel.ModelId)
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

func findAndUpdateJSON(ctx context.Context, modelDependencies map[string]any, key, value string, prismaClient *db.PrismaClient) (bool, error) {
	for k, v := range modelDependencies {
		if k == key {
			if val, ok := v.(string); ok && val == value {
				return updateFromDatabase(ctx, modelDependencies, key, value, prismaClient)
			}
		}
	}
	return false, nil
}

func updateFromDatabase(ctx context.Context, data map[string]any, key, value string, prismaClient *db.PrismaClient) (bool, error) {
	tracer := otel.Tracer("prisma")
	dbCtx, dbSpan := tracer.Start(ctx, "Get model")
	defer dbSpan.End()

	model, err := prismaClient.Model.FindUnique(
		db.Model.ID.Equals(value),
	).Exec(dbCtx)
	if err != nil {
		dbSpan.SetAttributes(attribute.String("error", err.Error()))
		return false, err
	}

	var attributes map[string]models.Attribute
	err = json.Unmarshal([]byte(model.Attributes), &attributes)
	if err != nil {
		return false, fmt.Errorf("failed to parse attributes: %w", err)
	}

	dbSpan.AddEvent("successfully got model from the database")

	// Replace the existing value with the fetched, unmarshalled model
	m := models.NewModel(attributes)

	data[key] = m

	return true, nil
}
