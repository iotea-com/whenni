package things

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/iotea-com/iotea/libs/engine/dependencies/certificates"
	"github.com/iotea-com/iotea/libs/secrets"
)

const MqttBrokerThingCategory ThingCategory = "MQTT_BROKER"
const MqttClientThingCategory ThingCategory = "MQTT_CLIENT"

type MqttBroker struct {
	Host     string              `json:"host" validate:"required,hostname|ip"`
	Port     int                 `json:"port" validate:"min=1001,max=65535"`
	Protocol string              `json:"protocol" validate:"required,oneof=mqtt mqtts"`
	CaCert   certificates.CaCert `json:"caCertificate"`
}

func NewMqttBrokerFromAttributes(attributes map[string]any) (*MqttBroker, error) {
	// Unmarshal
	jsonAttributes, _ := json.Marshal(attributes)
	var mqttBroker MqttBroker
	err := json.Unmarshal(jsonAttributes, &mqttBroker)
	if err != nil {
		return nil, err
	}

	// Validate
	err = mqttBroker.Validate()
	if err != nil {
		return nil, err
	}

	return &mqttBroker, nil
}

func (m *MqttBroker) Category() ThingCategory {
	return MqttBrokerThingCategory
}

func (m *MqttBroker) Validate() error {
	v := validator.New()

	if err := v.Struct(m); err != nil {
		return err
	}

	return nil
}

func (m *MqttBroker) MarshalJson() ([]byte, error) {
	jsonAttributes, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return jsonAttributes, nil
}

type MqttClient struct {
	Broker        any                        `json:"broker" validate:"required"` // Should reference a Thing UUID or be a Thing attributes struct itself
	Username      string                     `json:"username"`
	Password      string                     `json:"password"`
	ClientId      string                     `json:"clientId"`
	CertificateId string                     `json:"certificateId"`
	Keypair       certificates.ClientKeypair `json:"keypair,omitempty"` // This field is only used in the rules engine
}

func NewMqttClientFromAttributes(attributes map[string]any, secretsClient secrets.SecretsClient, spaceId string) (*MqttClient, error) {
	// Unmarshal
	jsonAttributes, _ := json.Marshal(attributes)
	var mqttClient MqttClient
	err := json.Unmarshal(jsonAttributes, &mqttClient)
	if err != nil {
		return nil, err
	}

	// Get password from secrets engine
	if mqttClient.Password != "" && mqttClient.Password != "__IOTEA_IGNORE__" {
		spaceSecrets, err := secretsClient.ReadKeyValue(context.Background(), fmt.Sprintf("spaces/%s", spaceId))
		if err != nil {
			return nil, fmt.Errorf("failed to resolve secrets: %v", err)
		}

		// Extract the secret from the space secrets
		mqttClient.Password = secretsClient.GetStringFromMap(spaceSecrets, mqttClient.Password)
	}

	// Validate
	err = mqttClient.Validate()
	if err != nil {
		return nil, err
	}

	return &mqttClient, nil
}

func (m *MqttClient) Category() ThingCategory {
	return MqttClientThingCategory
}

func (m *MqttClient) Validate() error {
	v := validator.New()

	if err := v.Struct(m); err != nil {
		return err
	}

	return nil
}

func (m *MqttClient) MarshalJson() ([]byte, error) {
	jsonAttributes, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return jsonAttributes, nil
}
