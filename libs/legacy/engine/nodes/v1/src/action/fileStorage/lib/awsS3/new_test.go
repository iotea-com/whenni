package awsS3

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

		// Assert the type of 'n' to be *S3ActionNode
		s3Node, ok := n.(*S3ActionSubnode)
		if !ok {
			t.Fatalf("New() returned incorrect type, expected *S3ActionNode")
		}

		// Assert that input channels are initialized correctly
		if len(s3Node.inputChannels) != 1 || s3Node.inputChannels[0].Type[0] != node.BytesDataType {
			t.Fatal("Input channels not initialized correctly")
		}
	})
}
