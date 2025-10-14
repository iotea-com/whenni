package main

import (
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	"github.com/iotea-com/iotea/libs/engine/nodes/v1/src/conditional/boolean/config"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

const (
	NodeName = "Boolean"
)

type BooleanConditionalNode struct {
	config         booleanConditionalNodeConfig.BooleanConditionalNodeConfig
	inputChannels  []node.IoChannel
	outputChannels []node.IoChannel
	notifyChannel  chan node.Notification

	// Observability fields
	tracer trace.Tracer
	meter  metric.Meter
	logger zerolog.Logger
}

func New() node.Interface {
	return &BooleanConditionalNode{
		inputChannels: []node.IoChannel{
			{
				Id:      "booleanInput",
				Channel: make(chan node.IoData, 100),
				Type: []node.DataType{
					node.BytesDataType,
					node.JsonDataType,
				},
			},
		},
		outputChannels: []node.IoChannel{
			{
				Id:      "Fail",
				Channel: make(chan node.IoData, 100),
				Type: []node.DataType{
					node.BytesDataType,
					node.JsonDataType,
				},
			},
			{
				Id:      "Pass",
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
