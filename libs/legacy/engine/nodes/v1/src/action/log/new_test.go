package main

import (
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

		// Assert the type of 'n' to be *LogActionNode
		logNode, ok := n.(*LogActionNode)
		if !ok {
			t.Errorf("New() returned incorrect type, expected *HttpActionNode")
		}

		// Assert that input channels are initialized correctly
		if len(logNode.inputChannels) != 1 || logNode.inputChannels[0].Type[0] != node.BytesDataType {
			t.Error("Input channels not initialized correctly")
		}

	})
}
