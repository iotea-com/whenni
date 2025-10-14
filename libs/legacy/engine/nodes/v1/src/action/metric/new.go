package main

import (
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"

	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	metricActionNodeConfig "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/action/metric/config"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

const (
	NodeName = "Metric"

	// TODO: Should this be a config option?
	CustomMetricsTable = "custom_metrics"
)

type MetricActionNode struct {
	config        metricActionNodeConfig.MetricActionNodeConfig
	inputChannels []node.IoChannel
	notifyChannel chan node.Notification

	// Clickhouse connection
	clickhouseClient *driver.Conn

	// Observability
	tracer trace.Tracer
	meter  metric.Meter
	logger zerolog.Logger
}

func New() node.Interface {
	return &MetricActionNode{
		inputChannels: []node.IoChannel{
			{
				Id:      "metricInput",
				Channel: make(chan node.IoData, 100),
				Type:    []node.DataType{node.BytesDataType},
			},
		},
	}
}

func main() {} // main is required by go, but node is built as plugin
