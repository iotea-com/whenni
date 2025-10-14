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

		// Assert the type of 'n' to be *TimerSourceNode
		timerNode, ok := n.(*TimerSourceNode)
		if !ok {
			t.Errorf("New() returned incorrect type, expected *TimerSourceNode")
		}

		// Assert that input channels are initialized correctly
		if len(timerNode.inputChannels) != 1 || timerNode.inputChannels[0].Type[0] != node.SignalDataType {
			t.Error("Input channels not initialized correctly")
		}

		// Assert that output channels are initialized correctly
		if len(timerNode.outputChannels) != 1 || timerNode.outputChannels[0].Type[0] != node.SignalDataType {
			t.Error("Output channels not initialized correctly")
		}
	})
}
