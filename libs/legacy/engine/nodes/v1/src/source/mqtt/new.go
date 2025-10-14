package main

import (
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/iotea-com/iotea/libs/engine/dependencies/things"
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	mqttSourceNodeConfig "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/source/mqtt/config"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

const (
	NodeName = "MQTT Subscriber"
)

// Remove these bars
const (
	SubscriptionTimeout = 1 * time.Second
	DisconnectTimeout   = 1 * time.Millisecond
)

type MqttSourceNodeMetrics struct {
	MessageCount     int
	MessageSizeBytes int
}

type MqttSourceNode struct {
	config         mqttSourceNodeConfig.MqttSourceNodeConfig
	inputChannels  []node.IoChannel
	outputChannels []node.IoChannel
	errChan        chan error
	client         mqtt.Client
	broker         things.MqttBroker
	notifyChannel  chan node.Notification

	// Observability
	tracer trace.Tracer
	meter  metric.Meter
	logger zerolog.Logger
}

func New() node.Interface {
	return &MqttSourceNode{
		inputChannels: []node.IoChannel{
			{
				Id:      "passThrough",
				Channel: make(chan node.IoData, 100),
				Type: []node.DataType{
					node.BytesDataType,
					node.JsonDataType,
				},
			},
		},
		outputChannels: []node.IoChannel{
			{
				Id:      "mqttOutput",
				Channel: make(chan node.IoData, 100),
				Type: []node.DataType{
					node.BytesDataType,
					node.JsonDataType,
				},
			},
		},
	}
}

func main() {} // main is required by go, but node is built as plugin
