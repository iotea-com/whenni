package kafkaSourceNode

import (
	"context"

	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
)

func (n *KafkaSourceSubnode) Deinit(params node.DeinitParams) node.Error {
	_, deinitSpan := n.tracer.Start(context.Background(), "Deinitialize Node")
	defer deinitSpan.End()

	// Close consumer
	if n.consumer != nil {
		n.consumer.Close()
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
