package main

import (
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	"github.com/rs/zerolog"
)

const (
	NodeName = "Time Series Database Action"
)

type TimeSeriesDbActionNode struct {
	inputChannels []node.IoChannel

	subnode node.Interface

	// Observability
	logger zerolog.Logger
}

func New() node.Interface {
	return &TimeSeriesDbActionNode{
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

func main() {} // main is required by go, but node is built as plugin
