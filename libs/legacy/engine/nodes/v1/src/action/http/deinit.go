package main

import (
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
)

func (n *HttpActionNode) Deinit(params node.DeinitParams) node.Error {
	n.client.CloseIdleConnections()

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
	return node.Error{
		Type: node.NoError,
	}
}
