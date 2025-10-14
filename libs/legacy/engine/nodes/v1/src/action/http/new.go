package main

import (
	"net/http"

	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	httpActionNodeConfig "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/action/http/config"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

const (
	NodeName = "HTTP Client"
)

type HttpActionNode struct {
	config         httpActionNodeConfig.HttpActionNodeConfig
	inputChannels  []node.IoChannel
	outputChannels []node.IoChannel
	notifyChannel  chan node.Notification
	client         http.Client

	// Observability
	tracer trace.Tracer
	meter  metric.Meter
	logger zerolog.Logger
}

func New() node.Interface {
	return &HttpActionNode{
		client: *http.DefaultClient,
		inputChannels: []node.IoChannel{
			{
				Id:      "httpInput",
				Channel: make(chan node.IoData, 100),
				Type:    []node.DataType{node.BytesDataType},
			},
		},
		outputChannels: []node.IoChannel{
			{
				Id:      "httpOutput",
				Channel: make(chan node.IoData, 100),
				Type:    []node.DataType{node.BytesDataType},
			},
		},
	}
}

func main() {} // main is required by go, but node is built as plugin
