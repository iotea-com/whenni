package channels

import (
	"slices"
	"testing"
)

func TestValidate(t *testing.T) {
	t.Run("Valid channel with one source node, multiple branches, and no hanging nodes", func(t *testing.T) {
		// Given
		channel := Channel{
			Nodes: []Node{
				{Id: "1", Metadata: NodeMetadata{Type: "source"}},
				{Id: "2", Metadata: NodeMetadata{Type: "process"}},
				{Id: "3", Metadata: NodeMetadata{Type: "action"}},
			},
			Edges: []Edge{
				{Id: "1", From: EdgeConnection{NodeId: "1", IoId: "mqttOutput"}, To: EdgeConnection{NodeId: "2", IoId: "bytesInput"}},
				{Id: "2", From: EdgeConnection{NodeId: "2", IoId: "bytesOutput"}, To: EdgeConnection{NodeId: "3", IoId: "bytesInput"}},
			},
		}

		// When
		validationErrors := channel.Validate()

		// Then
		if validationErrors != nil {
			t.Errorf("expected no error but got %v", validationErrors)
		}
	})

	t.Run("Invalid channel with no nodes", func(t *testing.T) {
		// Given
		channel := Channel{
			Nodes: []Node{},
			Edges: []Edge{},
		}

		// When
		validationErrors := channel.Validate()

		// Then
		if validationErrors == nil {
			t.Errorf("expected error but got nil")
		}

		expectedError := "Channel must have more than one node."
		if !slices.Contains(validationErrors.ChannelErrors, expectedError) {
			t.Errorf("expected error message %s but got %v", expectedError, validationErrors)
		}
	})

	t.Run("Invalid channel with only one node", func(t *testing.T) {
		// Given
		channel := Channel{
			Nodes: []Node{
				{Id: "1", Metadata: NodeMetadata{Type: "source"}},
			},
			Edges: []Edge{},
		}

		// When
		validationErrors := channel.Validate()

		// Then
		if validationErrors == nil {
			t.Errorf("expected error but got nil")
		}

		expectedError := "Channel must have more than one node."
		if !slices.Contains(validationErrors.ChannelErrors, expectedError) {
			t.Errorf("expected error message %s but got %v", expectedError, validationErrors)
		}
	})

	t.Run("Invalid channel with multiple source nodes", func(t *testing.T) {
		// Given
		channel := Channel{
			Nodes: []Node{
				{Id: "1", Metadata: NodeMetadata{Type: "source"}},
				{Id: "2", Metadata: NodeMetadata{Type: "source"}},
				{Id: "3", Metadata: NodeMetadata{Type: "action"}},
			},
			Edges: []Edge{
				{Id: "1", From: EdgeConnection{NodeId: "1", IoId: "bytesOutput"}, To: EdgeConnection{NodeId: "3", IoId: "bytesInput"}},
				{Id: "2", From: EdgeConnection{NodeId: "2", IoId: "bytesOutput"}, To: EdgeConnection{NodeId: "3", IoId: "bytesInput"}},
			},
		}

		// When
		validationErrors := channel.Validate()

		// Then
		if validationErrors == nil {
			t.Errorf("expected error but got nil")
		}

		expectedError := "Only one source node is allowed."
		if !slices.Contains(validationErrors.ChannelErrors, expectedError) {
			t.Errorf("expected error message %s but got %v", expectedError, validationErrors)
		}
	})

	t.Run("Invalid channel with a node having multiple inputs", func(t *testing.T) {
		// Given
		channel := Channel{
			Nodes: []Node{
				{Id: "1", Metadata: NodeMetadata{Type: "source", Name: "Fake Name"}},
				{Id: "2", Metadata: NodeMetadata{Type: "process", Name: "Fake Name"}},
				{Id: "3", Metadata: NodeMetadata{Type: "process", Name: "Fake Name"}},
				{Id: "4", Metadata: NodeMetadata{Type: "action", Name: "Fake Name"}},
			},
			Edges: []Edge{
				{Id: "1", From: EdgeConnection{NodeId: "1", IoId: "bytesOutput"}, To: EdgeConnection{NodeId: "2", IoId: "bytesInput"}},
				{Id: "2", From: EdgeConnection{NodeId: "2", IoId: "bytesOutput"}, To: EdgeConnection{NodeId: "3", IoId: "bytesInput"}},
				{Id: "3", From: EdgeConnection{NodeId: "2", IoId: "bytesOutput"}, To: EdgeConnection{NodeId: "4", IoId: "bytesInput"}},
				{Id: "4", From: EdgeConnection{NodeId: "3", IoId: "bytesOutput"}, To: EdgeConnection{NodeId: "4", IoId: "bytesInput"}},
			},
		}

		// When
		validationErrors := channel.Validate()

		// Then
		if validationErrors == nil {
			t.Errorf("expected error but got nil")
		}

		expectedError := "Node 4 (Fake Name) has more than one input."
		if !slices.Contains(validationErrors.ChannelErrors, expectedError) {
			t.Errorf("expected error message %s but got %v", expectedError, validationErrors)
		}
	})

	t.Run("Invalid channel with a node having no edges (hanging node)", func(t *testing.T) {
		// Given
		channel := Channel{
			Nodes: []Node{
				{Id: "1", Metadata: NodeMetadata{Type: "source", Name: "Fake Name"}},
				{Id: "2", Metadata: NodeMetadata{Type: "process", Name: "Fake Name"}},
				{Id: "3", Metadata: NodeMetadata{Type: "action", Name: "Fake Name"}},
				{Id: "4", Metadata: NodeMetadata{Type: "process", Name: "Fake Name"}},
			},
			Edges: []Edge{
				{Id: "1", From: EdgeConnection{NodeId: "1", IoId: "mqttOutput"}, To: EdgeConnection{NodeId: "2", IoId: "bytesInput"}},
				{Id: "2", From: EdgeConnection{NodeId: "2", IoId: "bytesOutput"}, To: EdgeConnection{NodeId: "3", IoId: "bytesInput"}},
			},
		}

		// When
		validationErrors := channel.Validate()

		// Then
		if validationErrors == nil {
			t.Errorf("expected error but got nil")
		}

		expectedError := "Node 4 (Fake Name) is a hanging node with no edges."

		if !slices.Contains(validationErrors.ChannelErrors, expectedError) {
			t.Errorf("expected error message %s but got %v", expectedError, validationErrors)
		}
	})

	t.Run("Invalid channel with an edge to a non-existent node", func(t *testing.T) {
		// Given
		channel := Channel{
			Nodes: []Node{
				{Id: "1", Metadata: NodeMetadata{Type: "source"}},
				{Id: "2", Metadata: NodeMetadata{Type: "process"}},
			},
			Edges: []Edge{
				{Id: "1", From: EdgeConnection{NodeId: "1", IoId: "mqttOutput"}, To: EdgeConnection{NodeId: "3", IoId: "bytesInput"}},
			},
		}

		// When
		validationErrors := channel.Validate()

		// Then
		if validationErrors == nil {
			t.Errorf("expected error but got nil")
		}

		expectedError := "Edge to node 3 does not exist."
		if !slices.Contains(validationErrors.ChannelErrors, expectedError) {
			t.Errorf("expected error message %s but got %v", expectedError, validationErrors)
		}
	})

	t.Run("Invalid channel with an edge from a non-existent node", func(t *testing.T) {
		// Given
		channel := Channel{
			Nodes: []Node{
				{Id: "1", Metadata: NodeMetadata{Type: "source"}},
				{Id: "2", Metadata: NodeMetadata{Type: "process"}},
			},
			Edges: []Edge{
				{Id: "1", From: EdgeConnection{NodeId: "3", IoId: "mqttOutput"}, To: EdgeConnection{NodeId: "2", IoId: "bytesInput"}},
			},
		}

		// When
		validationErrors := channel.Validate()

		// Then
		if validationErrors == nil {
			t.Errorf("expected error but got nil")
		}

		expectedError := "Edge from node 3 does not exist."

		if !slices.Contains(validationErrors.ChannelErrors, expectedError) {
			t.Errorf("expected error message %s but got %v", expectedError, validationErrors)
		}
	})

	t.Run("Invalid channel with an action node that has no input but has outputs", func(t *testing.T) {
		// Given
		channel := Channel{
			Nodes: []Node{
				{Id: "1", Metadata: NodeMetadata{Type: "source", Name: "Fake Name"}},
				{Id: "2", Metadata: NodeMetadata{Type: "process", Name: "Fake Name"}},
				{Id: "3", Metadata: NodeMetadata{Type: "action", Name: "Fake Name"}},
				{Id: "4", Metadata: NodeMetadata{Type: "process", Name: "Fake Name"}},
				{Id: "5", Metadata: NodeMetadata{Type: "action", Name: "Fake Name"}},
			},
			Edges: []Edge{
				{Id: "1", From: EdgeConnection{NodeId: "1", IoId: "mqttOutput"}, To: EdgeConnection{NodeId: "2", IoId: "bytesInput"}},
				{Id: "2", From: EdgeConnection{NodeId: "2", IoId: "bytesOutput"}, To: EdgeConnection{NodeId: "3", IoId: "bytesInput"}},
				{Id: "3", From: EdgeConnection{NodeId: "4", IoId: "bytesOutput"}, To: EdgeConnection{NodeId: "5", IoId: "bytesInput"}},
			},
		}

		// When
		validationErrors := channel.Validate()

		// Then
		if validationErrors == nil {
			t.Errorf("expected error but got nil")
		}

		expectedError := "Node 4 (Fake Name) has no input and is not a source node."
		if !slices.Contains(validationErrors.ChannelErrors, expectedError) {
			t.Errorf("expected error message %s but got %v", expectedError, validationErrors)
		}
	})

	t.Run("Invalid channel with an edge connecting two outputs", func(t *testing.T) {
		// Given
		channel := Channel{
			Nodes: []Node{
				{Id: "1", Metadata: NodeMetadata{Type: "source", Name: "Fake Name"}},
				{Id: "2", Metadata: NodeMetadata{Type: "process", Name: "Fake Name"}},
			},
			Edges: []Edge{
				{Id: "1", From: EdgeConnection{NodeId: "1", IoId: "mqttOutput"}, To: EdgeConnection{NodeId: "2", IoId: "bytesOutput"}},
			},
		}

		// When
		validationErrors := channel.Validate()

		// Then
		if validationErrors == nil {
			t.Errorf("expected error but got nil")
		}

		expectedError := "Edge 1 (between Fake Name and Fake Name) connects two outputs."
		if !slices.Contains(validationErrors.ChannelErrors, expectedError) {
			t.Errorf("expected error message %s but got %v", expectedError, validationErrors)
		}
	})

	t.Run("Invalid channel with an edge connecting two inputs", func(t *testing.T) {
		// Given
		channel := Channel{
			Nodes: []Node{
				{Id: "1", Metadata: NodeMetadata{Type: "source", Name: "Fake Name"}},
				{Id: "2", Metadata: NodeMetadata{Type: "process", Name: "Fake Name"}},
			},
			Edges: []Edge{
				{Id: "1", From: EdgeConnection{NodeId: "1", IoId: "mqttInput"}, To: EdgeConnection{NodeId: "2", IoId: "bytesInput"}},
			},
		}

		// When
		validationErrors := channel.Validate()

		// Then
		if validationErrors == nil {
			t.Errorf("expected error but got nil")
		}

		expectedError := "Edge 1 (between Fake Name and Fake Name) connects two inputs."
		if !slices.Contains(validationErrors.ChannelErrors, expectedError) {
			t.Errorf("expected error message %s but got %v", expectedError, validationErrors)
		}
	})

	t.Run("Invalid channel with a duplicate edge", func(t *testing.T) {
		// Given
		channel := Channel{
			Nodes: []Node{
				{Id: "1", Metadata: NodeMetadata{Type: "source"}},
				{Id: "2", Metadata: NodeMetadata{Type: "process"}},
			},
			Edges: []Edge{
				{Id: "1", From: EdgeConnection{NodeId: "1", IoId: "mqttInput"}, To: EdgeConnection{NodeId: "2", IoId: "bytesInput"}},
				{Id: "2", From: EdgeConnection{NodeId: "1", IoId: "mqttInput"}, To: EdgeConnection{NodeId: "2", IoId: "bytesInput"}},
			},
		}

		// When
		validationErrors := channel.Validate()

		// Then
		if validationErrors == nil {
			t.Errorf("expected error but got nil")
		}

		expectedError := "Duplicate edge detected between nodes 1 and 2 with IO points mqttInput and bytesInput."
		if !slices.Contains(validationErrors.ChannelErrors, expectedError) {
			t.Errorf("expected error message %s but got %v", expectedError, validationErrors)
		}
	})

	t.Run("Invalid channel with a same-node edge", func(t *testing.T) {
		// Given
		channel := Channel{
			Nodes: []Node{
				{Id: "1", Metadata: NodeMetadata{Type: "source", Name: "Fake Name"}},
				{Id: "2", Metadata: NodeMetadata{Type: "process", Name: "Fake Name"}},
			},
			Edges: []Edge{
				{Id: "1", From: EdgeConnection{NodeId: "1", IoId: "mqttInput"}, To: EdgeConnection{NodeId: "2", IoId: "bytesInput"}},
				{Id: "2", From: EdgeConnection{NodeId: "2", IoId: "bytesOutput"}, To: EdgeConnection{NodeId: "2", IoId: "bytesInput"}},
			},
		}

		// When
		validationErrors := channel.Validate()

		// Then
		if validationErrors == nil {
			t.Errorf("expected error but got nil")
		}

		expectedError := "Edge 2 (between Fake Name and Fake Name) connects an input to an output from the same node 2."
		if !slices.Contains(validationErrors.ChannelErrors, expectedError) {
			t.Errorf("expected error message %s but got %v", expectedError, validationErrors)
		}
	})
}
