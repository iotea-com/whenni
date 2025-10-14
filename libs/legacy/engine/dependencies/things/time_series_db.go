package things

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/iotea-com/iotea/libs/secrets"
)

const InfluxDbDatabaseThingCategory ThingCategory = "INFLUXDB_DATABASE"

type InfluxDbDatabase struct {
	Host     string `json:"host" validate:"required,hostname|ip"`
	Port     int    `json:"port" validate:"required"`
	Protocol string `json:"protocol" validate:"required,oneof=http https"`
	Token    string `json:"token" validate:"required"`
	OrgName  string `json:"org" validate:"required"`
}

func NewInfluxDbDatabaseFromAttributes(attributes map[string]any, secretsClient secrets.SecretsClient, spaceId string) (*InfluxDbDatabase, error) {
	// Unmarshal
	jsonAttributes, _ := json.Marshal(attributes)
	var influxDbDatabase InfluxDbDatabase
	err := json.Unmarshal(jsonAttributes, &influxDbDatabase)
	if err != nil {
		return nil, err
	}

	// Get password from secrets engine
	spaceSecrets, err := secretsClient.ReadKeyValue(context.Background(), fmt.Sprintf("spaces/%s", spaceId))
	if err != nil {
		return nil, fmt.Errorf("failed to resolve secrets: %v", err)
	}

	// Extract the secret from the space secrets
	influxDbDatabase.Token = secretsClient.GetStringFromMap(spaceSecrets, influxDbDatabase.Token)

	// Validate
	err = influxDbDatabase.Validate()
	if err != nil {
		return nil, err
	}

	return &influxDbDatabase, nil
}

func (m *InfluxDbDatabase) Category() ThingCategory {
	return InfluxDbDatabaseThingCategory
}

func (m *InfluxDbDatabase) Validate() error {
	v := validator.New()

	if err := v.Struct(m); err != nil {
		return err
	}

	return nil
}

func (m *InfluxDbDatabase) MarshalJson() ([]byte, error) {
	jsonAttributes, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return jsonAttributes, nil
}

const ClickhouseDatabaseThingCategory ThingCategory = "CLICKHOUSE_DATABASE"

type ClickhouseDatabase struct {
	Host     string `json:"host" validate:"required,hostname|ip"`
	Port     int    `json:"port" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Database string `json:"database" validate:"required"`
}

func NewClickhouseDatabaseFromAttributes(attributes map[string]any, secretsClient secrets.SecretsClient, spaceId string) (*ClickhouseDatabase, error) {
	// Unmarshal
	jsonAttributes, _ := json.Marshal(attributes)
	var clickhouseDatabase ClickhouseDatabase
	err := json.Unmarshal(jsonAttributes, &clickhouseDatabase)
	if err != nil {
		return nil, err
	}

	// Get password from secrets engine
	spaceSecrets, err := secretsClient.ReadKeyValue(context.Background(), fmt.Sprintf("spaces/%s", spaceId))
	if err != nil {
		return nil, fmt.Errorf("failed to resolve secrets: %v", err)
	}

	// Extract the secret from the space secrets
	clickhouseDatabase.Password = secretsClient.GetStringFromMap(spaceSecrets, clickhouseDatabase.Password)

	// Validate
	err = clickhouseDatabase.Validate()
	if err != nil {
		return nil, err
	}

	return &clickhouseDatabase, nil
}

func (m *ClickhouseDatabase) Category() ThingCategory {
	return ClickhouseDatabaseThingCategory
}

func (m *ClickhouseDatabase) Validate() error {
	v := validator.New()

	if err := v.Struct(m); err != nil {
		return err
	}

	return nil
}

func (m *ClickhouseDatabase) MarshalJson() ([]byte, error) {
	jsonAttributes, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return jsonAttributes, nil
}
