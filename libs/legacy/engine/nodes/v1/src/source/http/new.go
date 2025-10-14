package main

import (
	"github.com/gofiber/fiber/v2"
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	httpSourceNodeConfig "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/source/http/config"
	router "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/source/http/lib"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

const (
	NodeName = "HTTP Server"

	// Channel IDs
	PassThroughChannelId    = "passThrough"
	ResponsesInputChannelId = "httpResponseInput"
	HttpOutputChannelId     = "httpOutput"
)

type HttpSourceNode struct {
	config         httpSourceNodeConfig.HttpSourceNodeConfig
	inputChannels  []node.IoChannel
	outputChannels []node.IoChannel
	errChan        chan error
	server         *fiber.App
	notifyChannel  chan node.Notification

	// Observability
	tracer trace.Tracer
	meter  metric.Meter
	logger zerolog.Logger

	router *router.ResponseRouter
}

func New() node.Interface {
	return &HttpSourceNode{
		inputChannels: []node.IoChannel{
			{
				Id:      PassThroughChannelId,
				Channel: make(chan node.IoData, 100),
				Type: []node.DataType{
					node.BytesDataType,
					node.JsonDataType,
					node.FileDataType,
				},
			},
			{
				Id:      ResponsesInputChannelId,
				Channel: make(chan node.IoData, 100),
				Type: []node.DataType{
					node.BytesDataType,
				},
			},
		},
		outputChannels: []node.IoChannel{
			{
				Id:      HttpOutputChannelId,
				Channel: make(chan node.IoData, 100),
				Type: []node.DataType{
					node.BytesDataType,
					node.JsonDataType,
					node.FileDataType,
				},
			},
		},
	}
}

func main() {} // main is required by go, but node is built as plugin
