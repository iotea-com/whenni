package minio

import (
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
)

func (n *MinioActionSubnode) Deinit(params node.DeinitParams) node.Error {

	// Close nodes safely
	safeClose := func(ch chan node.IoData) {
		defer func() {
			if r := recover(); r != nil {
				// "Attempted to close an already-closed node"
			}
		}()
		close(ch)
	}

	// Close all input nodes
	for _, ch := range n.inputChannels {
		safeClose(ch.Channel)
	}

	// No output nodes to close
	return node.Error{
		Type: node.NoError,
	}
}
