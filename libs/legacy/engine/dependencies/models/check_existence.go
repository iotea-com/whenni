package models

import (
	"encoding/json"
	"fmt"
)

func (m *Model) CheckExistence(input map[string]json.RawMessage, attributeId string) error {
	// Validate that the input data matches the input struct
	if err := m.Compare(input); err != nil {
		// Marshal the struct input
		inputSchemaBytes, jsonErr := json.Marshal(m)
		if jsonErr != nil {
			return fmt.Errorf("failed to marshal input schema during comparison error: %s (original error: %s)", jsonErr, err)
		}

		inputBytes, jsonErr := json.Marshal(input)
		if jsonErr != nil {
			return fmt.Errorf("failed to marshal input data during comparison error: %s (original error: %s)", jsonErr, err)
		}

		return fmt.Errorf("input data validation failed: %v, expected struct format: %v, received: %v", err, string(inputSchemaBytes), string(inputBytes))
	}

	m.ensureCaches() // Ensure caches are initialized

	// Build the path to the attribute if not already cached
	if _, exists := m.attributePaths[attributeId]; !exists {
		m.buildAttributePath(attributeId)
	}

	// Get the path for this attribute
	path := m.attributePaths[attributeId]
	if len(path) == 0 {
		return fmt.Errorf("attribute \"%s\" not found in model", attributeId)
	}

	// Get the nested value using the cached path
	_, exists := m.getNestedValue(input, path)
	if !exists {
		return fmt.Errorf("attribute \"%s\" (%s) not found in input", attributeId, path)
	}

	return nil
}
