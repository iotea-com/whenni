package main

import (
	"testing"
	"time"
)

const (
	ShortLifetimeExpectedDuration = 100 * time.Millisecond
)

// TestInit_Config tests that the node won't accept a wrong configuration byte array passed
func TestInit_Config(t *testing.T) {
	t.Run("returns error on wrong config", func(t *testing.T) {
		// TODO: Implement this test
	})
}

// TestInit tests that the action node will NOT connect to the broker and
// return immediately
func TestInit(t *testing.T) {
	t.Parallel()

	t.Run("successfully initializes node", func(t *testing.T) {
		// TODO: Implement this test
	})
}
