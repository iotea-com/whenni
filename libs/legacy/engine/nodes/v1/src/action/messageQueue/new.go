package main

import (
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	"github.com/rs/zerolog"
)

const (
	NodeName = "Message Queue Action"
)

type MessageQueueActionNode struct {
	inputChannels []node.IoChannel

	subnode node.Interface

	// Observability
	logger zerolog.Logger
}

func New() node.Interface {
	return &MessageQueueActionNode{
		inputChannels: []node.IoChannel{
			{
				Id:      "messageQueueInput",
				Channel: make(chan node.IoData, 100),
				Type:    []node.DataType{node.BytesDataType},
			},
		},
	}
}

func main() {} // main is required by go, but node is built as plugin
