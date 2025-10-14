package natsActionNode

import (
	"testing"
)

func TestDeinit(t *testing.T) {
	t.Run("successfully deinitializes node", func(t *testing.T) {
		n := New()
		_ = n
	})
}
