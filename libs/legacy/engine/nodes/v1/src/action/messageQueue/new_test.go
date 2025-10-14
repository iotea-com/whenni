package main

import (
	"context"
	"testing"

	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
)

func TestNew(t *testing.T) {
	t.Run("successfully creates new message queue node", func(t *testing.T) {
		n := New()

		// Assert that 'n' is not nil
		if n == nil {
			t.Error("New() returned nil, expected non-nil node.Interface")
		}

		// Assert the type of 'n' to be *MessageQueueActionNode
		messageQueueNode, ok := n.(*MessageQueueActionNode)
		if !ok {
			t.Errorf("New() returned incorrect type, expected *MessageQueueActionNode")
		}

		// Assert that input channels are initialized correctly
		if len(messageQueueNode.GetInputChans(context.Background())) != 1 || messageQueueNode.GetInputChans(context.Background())[0].Type[0] != node.BytesDataType {
			t.Error("Input channels not initialized correctly")
		}

		// Assert that output channels are initialized correctly
		if len(messageQueueNode.GetOutputChans(context.Background())) != 0 {
			t.Error("Node should have no output channels")
		}
	})
}
