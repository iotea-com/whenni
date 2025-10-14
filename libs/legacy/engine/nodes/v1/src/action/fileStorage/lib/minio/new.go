package minio

import (
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	minioActionNodeConfig "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/action/fileStorage/lib/minio/config"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

const (
	NodeName = "Minio File Storage"
)

type MinioActionSubnode struct {
	config        minioActionNodeConfig.MinioActionSubnodeConfig
	inputChannels []node.IoChannel
	notifyChannel chan node.Notification

	// Observability fields
	tracer trace.Tracer
	meter  metric.Meter
	logger zerolog.Logger
}

func New() node.Interface {
	return &MinioActionSubnode{
		inputChannels: []node.IoChannel{
			{
				Id:      "fileStorageInput",
				Channel: make(chan node.IoData, 100),
				Type: []node.DataType{
					node.BytesDataType,
					node.JsonDataType,
					node.FileDataType,
				},
			},
		},
		// No output channels
	}
}

func main() {}
