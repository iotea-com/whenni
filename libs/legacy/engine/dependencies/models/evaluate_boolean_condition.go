package models

import (
	"encoding/json"
	"fmt"
)

type BooleanCondition struct {
	AttributeId string `json:"attributeId" validate:"required"`
	Operator    string `json:"operator" validate:"required,oneof=equal-to"`
	Value       bool   `json:"value"`
}

func (m *Model) EvaluateBooleanCondition(input map[string]json.RawMessage, condition BooleanCondition) error {
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
	if _, exists := m.attributePaths[condition.AttributeId]; !exists {
		m.buildAttributePath(condition.AttributeId)
	}

	// Get the path for this attribute
	path := m.attributePaths[condition.AttributeId]
	if len(path) == 0 {
		return fmt.Errorf("attribute \"%s\" not found in model", condition.AttributeId)
	}

	// Get the nested value using the cached path
	val, exists := m.getNestedValue(input, path)
	if !exists {
		return fmt.Errorf("attribute \"%s\" not found in input", condition.AttributeId)
	}

	// Convert the value to bool for comparison
	var boolVal bool
	switch v := val.(type) {
	case bool:
		boolVal = v
	default:
		return fmt.Errorf("attribute %s is not a bool", condition.AttributeId)
	}

	switch condition.Operator {
	case "equal-to":
		if boolVal == condition.Value {
			return nil // Condition met
		} else {
			return fmt.Errorf("attribute %s is not equal to %v", condition.AttributeId, condition.Value)
		}
	default:
		return fmt.Errorf("unsupported operator %s", condition.Operator)
	}
}
