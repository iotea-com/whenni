package main

import (
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/iotea-com/iotea/libs/engine/dependencies/things"
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	mqttActionNodeConfig "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/action/mqtt/config"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

const (
	NodeName = "MQTT Publisher"
)

// Remove these bars
const (
	// SubscriptionTimeout = 1 * time.Second
	DisconnectTimeout = 1 * time.Millisecond
)

type MqttActionNodeMetrics struct {
	MessageCount     int
	MessageSizeBytes int
}

type MqttActionNode struct {
	config        mqttActionNodeConfig.MqttActionNodeConfig
	inputChannels []node.IoChannel
	notifyChannel chan node.Notification

	client mqtt.Client
	broker things.MqttBroker

	// Observability
	tracer trace.Tracer
	meter  metric.Meter
	logger zerolog.Logger
}

func New() node.Interface {
	return &MqttActionNode{
		inputChannels: []node.IoChannel{
			{
				Id:      "mqttInput",
				Channel: make(chan node.IoData, 100),
				Type:    []node.DataType{node.BytesDataType},
			},
		},
	}
}

func main() {} // main is required by go, but node is built as plugin
