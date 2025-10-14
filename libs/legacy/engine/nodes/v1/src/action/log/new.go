package main

import (
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

const (
	NodeName = "Log"
)

type LogActionNode struct {
	inputChannels []node.IoChannel
	notifyChannel chan node.Notification

	// Observability
	tracer trace.Tracer
	meter  metric.Meter
	logger zerolog.Logger
}

func New() node.Interface {
	return &LogActionNode{
		inputChannels: []node.IoChannel{
			{
				Id:      "logInput",
				Channel: make(chan node.IoData, 100),
				Type:    []node.DataType{node.BytesDataType},
			},
		},
	}
}

func main() {} // main is required by go, but node is built as plugin
