package channels

import "encoding/json"

// Edges
type Edge struct {
	Id   string         `json:"id"`
	From EdgeConnection `json:"from"`
	To   EdgeConnection `json:"to"`
}

type EdgeConnection struct {
	NodeId string `json:"nodeId"`
	IoId   string `json:"ioId"`
}

// Node
type Node struct {
	Id          string          `json:"id"`
	Coordinates NodeCoordinates `json:"coordinates"`
	Metadata    NodeMetadata    `json:"metadata"`
}

type NodeCoordinates struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type NodeModelDependency struct {
	ModelId   string `json:"modelId"`
	FieldName string `json:"fieldName"` // Name of the `config` field where the model belongs
}

type NodeThingDependency struct {
	Internal  bool   `json:"internal"` // Flag to indicate whether this thing is displayed to the user
	ThingId   string `json:"thingId"`
	FieldName string `json:"fieldName"` // Name of the `config` field where the thing belongs
}

type NodeDependencies struct {
	Models []NodeModelDependency `json:"models"`
	Things []NodeThingDependency `json:"things"`
}

type NodeMetadata struct {
	Version      string           `json:"version"`
	Name         string           `json:"name"`
	Label        string           `json:"label"`
	Type         string           `json:"type"`
	Description  string           `json:"description"`
	Io           NodeIo           `json:"io"`
	Config       json.RawMessage  `json:"config"`
	Dependencies NodeDependencies `json:"dependencies"`
}

type NodeIo struct {
	Inputs  []Io `json:"inputs"`
	Outputs []Io `json:"outputs"`
}

type Io struct {
	Id           string   `json:"id"`
	AllowedNodes []string `json:"allowedEdge"`
	DataType     []string `json:"dataType"`
}

type Runtime struct {
	Size string `json:"size"` // small, medium, large
}

// Channel
type Channel struct {
	Id      string  `json:"id"`
	Name    string  `json:"name"`
	Runtime Runtime `json:"runtime"`
	Nodes   []Node  `json:"nodes"`
	Edges   []Edge  `json:"edges"`
}
