package main

import (
	"context"
	"errors"

	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
)

func (n *MessageQueueSourceNode) GetInputChans(ctx context.Context) []node.IoChannel {
	return n.inputChannels
}

func (n *MessageQueueSourceNode) SetInputChan(ctx context.Context, from *node.IoChannel, to *node.IoChannel) error {
	if !containsChannel(n.inputChannels, to) {
		return errors.New("the 'to' channel is not an input channel of this node")
	}
	to.Channel = from.Channel // Connect the 'to' channel to the 'from' channel
	return nil
}

func (n *MessageQueueSourceNode) GetOutputChans(ctx context.Context) []node.IoChannel {
	return n.outputChannels
}

func (n *MessageQueueSourceNode) SetOutputChan(ctx context.Context, to *node.IoChannel, from *node.IoChannel) error {
	if !containsChannel(n.outputChannels, to) {
		return errors.New("the 'to' channel is not an output channel of this node")
	}
	from.Channel = to.Channel // Direct data sent to the 'to' channel to the 'from' channel of the next node
	return nil
}

// containsChannel Checks if a channel is in a slice of channels
func containsChannel(slice []node.IoChannel, channel *node.IoChannel) bool {
	for _, ch := range slice {
		if ch.Id == channel.Id {
			return true
		}
	}
	return false
}
