package main

import (
	"context"
	"errors"

	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
)

func (n *TransformProcessingNode) GetInputChans(ctx context.Context) []node.IoChannel {
	return n.inputChannels
}

// Utility function to check if a channel is in a slice of channels
func containsChannel(slice []node.IoChannel, channel *node.IoChannel) bool {
	for _, ch := range slice {
		if ch.Id == channel.Id {
			return true
		}
	}
	return false
}

// SetInputChan connects an output channel of another node (from) to an input channel of this node (to).
func (n *TransformProcessingNode) SetInputChan(ctx context.Context, from *node.IoChannel, to *node.IoChannel) error {
	if !containsChannel(n.inputChannels, to) {
		return errors.New("the 'to' channel is not an input channel of this node")
	}
	to.Channel = from.Channel // Connect the 'to' channel to the 'from' channel
	return nil
}

// SetOutputChan connects an output channel of this node (to) to an input channel of another node (from).
func (n *TransformProcessingNode) SetOutputChan(ctx context.Context, to *node.IoChannel, from *node.IoChannel) error {
	if !containsChannel(n.outputChannels, to) {
		return errors.New("the 'to' channel is not an output channel of this node")
	}
	from.Channel = to.Channel // Direct data sent to the 'to' channel to the 'from' channel of the next node
	return nil
}

func (n *TransformProcessingNode) GetOutputChans(ctx context.Context) []node.IoChannel {
	return n.outputChannels
}
