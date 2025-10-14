package main

import (
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
)

func (n *DocumentDbActionNode) Deinit(params node.DeinitParams) node.Error {
	// Close nodes safely
	safeClose := func(ch chan node.IoData) {
		defer func() {
			recover()
		}()
		close(ch)
	}

	// Close all input nodes
	for _, ch := range n.inputChannels {
		safeClose(ch.Channel)
	}

	// Close all output nodes
	for _, ch := range n.outputChannels {
		safeClose(ch.Channel)
	}

	// Call deinit on subnode
	if n.subnode != nil {
		return n.subnode.Deinit(params)
	}

	return node.Error{
		Type: node.NoError,
	}
}
