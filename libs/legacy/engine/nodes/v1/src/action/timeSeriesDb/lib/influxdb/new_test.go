package influxdbNode

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
			t.Fatal("New() returned nil, expected non-nil node.Interface")
		}

		// Assert the type of 'n' to be *InfluxDBActionNode
		influxdbNode, ok := n.(*InfluxDBActionNode)
		if !ok {
			t.Fatalf("New() returned incorrect type, expected *InfluxDBActionNode")
		}

		// Assert that input channels are initialized correctly
		inputChans := influxdbNode.GetInputChans(context.Background())
		if len(inputChans) != 1 || !matchDataTypes(inputChans[0].Type, []node.DataType{node.BytesDataType, node.JsonDataType}) {
			t.Error("Input channels not initialized correctly")
		}

		// Assert that output channels are initialized correctly
		if len(influxdbNode.GetOutputChans(context.Background())) != 0 {
			t.Error("Node should have no output channels")
		}
	})
}

// Helper function to match data types
func matchDataTypes(a, b []node.DataType) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
