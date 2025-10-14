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

		// Assert the type of 'n' to be *MqttSourceNode
		mqttNode, ok := n.(*MqttSourceNode)
		if !ok {
			t.Errorf("New() returned incorrect type, expected *MqttSourceNode")
		}

		// Assert that input channels are initialized correctly
		if len(mqttNode.inputChannels) != 1 || mqttNode.inputChannels[0].Type[0] != node.BytesDataType {
			t.Error("Input channels not initialized correctly")
		}

		// Assert that output channels are initialized correctly
		if len(mqttNode.outputChannels) != 1 || mqttNode.outputChannels[0].Type[0] != node.BytesDataType {
			t.Error("Output channels not initialized correctly")
		}
	})
}
