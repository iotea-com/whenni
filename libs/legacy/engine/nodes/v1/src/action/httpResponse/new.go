package main

import (
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	httpResponseActionNodeConfig "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/action/httpResponse/config"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

const (
	NodeName = "HTTP Response"
)

type HttpResponseActionNode struct {
	config         httpResponseActionNodeConfig.HttpResponseActionNodeConfig
	inputChannels  []node.IoChannel
	outputChannels []node.IoChannel
	notifyChannel  chan node.Notification

	// Observability
	tracer trace.Tracer
	meter  metric.Meter
	logger zerolog.Logger
}

func New() node.Interface {
	return &HttpResponseActionNode{
		inputChannels: []node.IoChannel{
			{
				Id:      "httpResponseInput",
				Channel: make(chan node.IoData, 100),
				Type:    []node.DataType{node.BytesDataType},
			},
		},
		outputChannels: []node.IoChannel{
			{
				Id:      "httpResponseOutput",
				Channel: make(chan node.IoData, 100),
				Type:    []node.DataType{node.BytesDataType},
			},
		},
	}
}

func main() {} // main is required by go, but node is built as plugin
