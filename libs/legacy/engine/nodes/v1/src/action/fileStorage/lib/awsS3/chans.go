package awsS3

import (
	"context"
	"errors"

	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
)

func (n *S3ActionSubnode) GetInputChans(ctx context.Context) []node.IoChannel {
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
func (n *S3ActionSubnode) SetInputChan(ctx context.Context, from *node.IoChannel, to *node.IoChannel) error {
	if !containsChannel(n.inputChannels, to) {
		return errors.New("the 'to' channel is not an input channel of this node")
	}
	to.Channel = from.Channel // Connect the 'to' channel to the 'from' channel
	return nil
}

// s3 action node doesn't have output channel
func (n *S3ActionSubnode) SetOutputChan(ctx context.Context, to *node.IoChannel, from *node.IoChannel) error {
	return nil
}

func (n *S3ActionSubnode) GetOutputChans(ctx context.Context) []node.IoChannel {
	return nil
}
