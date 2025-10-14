package natsActionNode

import (
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
)

func (n *NatsActionSubnode) Deinit(params node.DeinitParams) node.Error {
	// Close client connection
	if n.conn != nil {
		n.conn.Close()
	}

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
