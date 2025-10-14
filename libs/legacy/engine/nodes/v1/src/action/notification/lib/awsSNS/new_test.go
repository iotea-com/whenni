package awsSNS

import (
	"context"
	"testing"

	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
)

func TestNew(t *testing.T) {
	t.Run("successfully creates new node", func(t *testing.T) {
		n := New()

		// Assert that 'n' is not nil
		if n == nil {
			t.Error("New() returned nil, expected non-nil node.Interface")
		}

		// Assert the type of 'n' to be *KafkaActionNode
		messageQueueNode, ok := n.(*AwsSNSActionSubnode)
		if !ok {
			t.Errorf("New() returned incorrect type, expected *HttpActionNode")
		}

		// Assert that input channels are initialized correctly
		if len(messageQueueNode.inputChannels) != 1 || messageQueueNode.inputChannels[0].Type[0] != node.BytesDataType {
			t.Error("Input channels not initialized correctly")
		}

		// Assert that input channels are initialized correctly
		if len(messageQueueNode.GetOutputChans(context.Background())) != 0 {
			t.Error("Node should have no output channels")
		}
	})
}
