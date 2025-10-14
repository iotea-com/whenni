package main

import (
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	"github.com/rs/zerolog"
)

const (
	NodeName = "Document Database Action"
)

type DocumentDbActionNode struct {
	inputChannels  []node.IoChannel
	outputChannels []node.IoChannel

	subnode node.Interface

	// Observability
	logger zerolog.Logger
}

func New() node.Interface {
	return &DocumentDbActionNode{
		inputChannels: []node.IoChannel{
			{
				Id:      "documentDbInput",
				Channel: make(chan node.IoData, 100),
				Type:    []node.DataType{node.BytesDataType},
			},
		},
		outputChannels: []node.IoChannel{
			{
				Id:      "documentDbOutput",
				Channel: make(chan node.IoData, 100),
				Type:    []node.DataType{node.BytesDataType},
			},
		},
	}
}

func main() {} // main is required by go, but node is built as plugin
