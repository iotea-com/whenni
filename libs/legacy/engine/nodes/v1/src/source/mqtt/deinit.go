package main

import (
	"fmt"

	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
)

func (n *MqttSourceNode) Deinit(params node.DeinitParams) node.Error {

	if n.client != nil && n.client.IsConnected() {
		// first we unsubscribe
		unSubtoken := n.client.Unsubscribe(n.config.Topic)

		if unSubtoken.Wait() && unSubtoken.Error() != nil {
			return node.Error{
				Type:   node.FatalError,
				Reason: fmt.Sprintf("error unsubscribing from topic %s: %v", n.config.Topic, unSubtoken.Error()),
			}
		}

		// then we disconnect
		n.client.Disconnect(uint(DisconnectTimeout))
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

	// Close all output nodes
	for _, ch := range n.outputChannels {
		safeClose(ch.Channel)
	}

	return node.Error{
		Type: node.NoError,
	}
}
