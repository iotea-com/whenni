package influxdbNode

import (
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/domain"
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	influxdbActionNodeConfig "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/action/timeSeriesDb/lib/influxdb/config"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

const (
	NodeName = "InfluxDB Time Series DB"
)

type InfluxDBActionNode struct {
	config        influxdbActionNodeConfig.InfluxDBActionNodeConfig
	inputChannels []node.IoChannel
	notifyChannel chan node.Notification

	client     influxdb2.Client
	writeAPI   api.WriteAPIBlocking
	bucketsAPI api.BucketsAPI
	orgsAPI    api.OrganizationsAPI
	org        *domain.Organization

	// Observability fields
	tracer trace.Tracer
	meter  metric.Meter
	logger zerolog.Logger
}

func New() node.Interface {
	return &InfluxDBActionNode{
		inputChannels: []node.IoChannel{
			{
				Id:      "timeSeriesDbInput",
				Channel: make(chan node.IoData, 100),
				Type: []node.DataType{
					node.BytesDataType,
					node.JsonDataType,
				},
			},
		},
	}
}

func main() {}
