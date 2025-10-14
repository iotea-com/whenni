package minio

import (
	"testing"

	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
)

func TestNew(t *testing.T) {
	t.Run("successfully creates new node", func(t *testing.T) {
		n := New()
		_ = n

		// Assert that 'n' is not nil
		if n == nil {
			t.Fatal("New() returned nil, expected non-nil node.Interface")
		}

		// Assert the type of 'n' to be *MinioActionNode
		minioNode, ok := n.(*MinioActionSubnode)
		if !ok {
			t.Fatalf("New() returned incorrect type, expected *MinioActionNode")
		}

		// Assert that input channels are initialized correctly
		if len(minioNode.inputChannels) != 1 || minioNode.inputChannels[0].Type[0] != node.BytesDataType {
			t.Fatal("Input channels not initialized correctly")
		}
	})
}
