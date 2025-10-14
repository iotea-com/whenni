package main

import (
	"testing"
)

// Define a struct representing your data
type TestData struct {
	ID string `json:"id"`
}

func TestExec_LongLifetime(t *testing.T) {
	t.Run("successfully executes node with long lifetime", func(t *testing.T) {
	})
}

func TestExec_FailsWithNonByteData(t *testing.T) {
	t.Run("successfully fails when non-byte data is passed", func(t *testing.T) {
	})
}
