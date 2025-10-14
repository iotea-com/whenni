package things

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/iotea-com/iotea/libs/secrets"
)

const MongoDbServerThingCategory ThingCategory = "MONGODB_SERVER"

type MongoDbServer struct {
	Protocol string `json:"protocol" validate:"required,oneof=mongodb mongodb+srv"`
	Host     string `json:"host" validate:"required,hostname|ip"`
	Port     int    `json:"port" validate:"required_if=Protocol mongodb"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func NewMongoDbServerFromAttributes(attributes map[string]any, secretsClient secrets.SecretsClient, spaceId string) (*MongoDbServer, error) {
	// Unmarshal
	jsonAttributes, _ := json.Marshal(attributes)
	var mongoDbServer MongoDbServer
	err := json.Unmarshal(jsonAttributes, &mongoDbServer)
	if err != nil {
		return nil, err
	}

	// Get password from secrets engine
	spaceSecrets, err := secretsClient.ReadKeyValue(context.Background(), fmt.Sprintf("spaces/%s", spaceId))
	if err != nil {
		return nil, fmt.Errorf("failed to resolve secrets: %v", err)
	}

	// Extract the secret from the space secrets
	mongoDbServer.Password = secretsClient.GetStringFromMap(spaceSecrets, mongoDbServer.Password)

	// Validate
	err = mongoDbServer.Validate()
	if err != nil {
		return nil, err
	}

	return &mongoDbServer, nil
}

func (m *MongoDbServer) Category() ThingCategory {
	return MongoDbServerThingCategory
}

func (m *MongoDbServer) Validate() error {
	v := validator.New()

	if err := v.Struct(m); err != nil {
		return err
	}

	return nil
}

func (m *MongoDbServer) MarshalJson() ([]byte, error) {
	jsonAttributes, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return jsonAttributes, nil
}
