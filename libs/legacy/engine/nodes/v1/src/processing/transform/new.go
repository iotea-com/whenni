package main

import (
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	transformProcessingNodeConfig "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/processing/transform/config"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

const (
	NodeName = "Data Transform"
)

type TransformProcessingNode struct {
	config         transformProcessingNodeConfig.TransformProcessingNodeConfig
	inputChannels  []node.IoChannel
	outputChannels []node.IoChannel
	notifyChannel  chan node.Notification

	// Observability fields
	tracer trace.Tracer
	meter  metric.Meter
	logger zerolog.Logger
}

func New() node.Interface {
	return &TransformProcessingNode{
		inputChannels: []node.IoChannel{
			{
				Id:      "transformInput",
				Channel: make(chan node.IoData, 100),
				Type: []node.DataType{
					node.BytesDataType,
					node.JsonDataType,
				},
			},
		},
		outputChannels: []node.IoChannel{
			{
				Id:      "transformOutput",
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
