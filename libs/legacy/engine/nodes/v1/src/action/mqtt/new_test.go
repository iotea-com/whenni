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
			t.Fatal("New() returned nil, expected non-nil node.Interface")
		}

		// Assert the type of 'n' to be *MqttActionNode
		mqttNode, ok := n.(*MqttActionNode)
		if !ok {
			t.Fatalf("New() returned incorrect type, expected *MqttActionNode")
		}

		// Assert that input channels are initialized correctly
		if len(mqttNode.inputChannels) != 1 || mqttNode.inputChannels[0].Type[0] != node.BytesDataType {
			t.Fatal("Input channels not initialized correctly")
		}

	})
}
