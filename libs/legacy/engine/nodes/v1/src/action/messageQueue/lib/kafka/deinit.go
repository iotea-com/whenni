package kafkaActionNode

import (
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
)

func (n *KafkaActionSubnode) Deinit(params node.DeinitParams) node.Error {
	// Close producer
	if n.producer != nil {
		n.producer.Close()
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
