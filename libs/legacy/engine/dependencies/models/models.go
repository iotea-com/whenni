package models

import (
	"encoding/json"
	"sync"
)

type Attribute struct {
	Id                    string               `json:"id" validate:"required"`
	Key                   string               `json:"key" validate:"required,isValidFieldKey"`
	Type                  string               `json:"type" validate:"required,oneof=string number boolean array object"`
	Subtype               string               `json:"subtype" validate:"required_if=Type array"`
	ArrayObjectAttributes map[string]Attribute `json:"arrayObjectAttributes" validate:"required_if=Subtype object"`
	ParentId              string               `json:"parentId"` // Present if this attribute is a child of another attribute
	Required              bool                 `json:"required"`
	Protected             bool                 `json:"protected"`
	Order                 int                  `json:"order"`
}

type Model struct {
	Attributes map[string]Attribute `json:"attributes"`

	// Caches for optimized TransformData
	attributePaths  map[string][]string                   // Cache for attribute paths
	parentObjects   map[string]map[string]json.RawMessage // Cache for unmarshaled parent objects
	nestedMaps      map[string]map[string]any             // Cache for nested output maps
	cacheBuildMutex sync.RWMutex                          // Add mutex for thread safety
}

/*
Initializes the schema and computes all path caching.

If attempting to modify the model after initialization, you must create a new Model instance
with the modified attributes.
*/
func NewModel(attributes map[string]Attribute) *Model {
	return &Model{Attributes: attributes}
}
