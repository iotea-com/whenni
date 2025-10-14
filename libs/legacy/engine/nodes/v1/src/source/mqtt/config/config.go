package mqttSourceNodeConfig

import "github.com/iotea-com/iotea/libs/engine/dependencies/things"

type MqttSourceNodeConfig struct {
	ThingMqttClient things.MqttClient `json:"mqttClient::thing" validate:"required"`

	// Subscription Information
	Topic string `json:"topic" validate:"required"`
	QoS   int    `json:"qos" validate:"min=0,max=2"`
}
