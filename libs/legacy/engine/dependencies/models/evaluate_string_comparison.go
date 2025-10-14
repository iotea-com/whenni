package models

import (
	"encoding/json"
	"fmt"
	"strings"
)

type StringComparison struct {
	AttributeId string `json:"attributeId" validate:"required"`
	Value       string `json:"value" validate:"required"`
	Operator    string `json:"operator" validate:"required,oneof=equal-to contains starts-with ends-with"`
}

func (m *Model) EvaluateStringComparison(input map[string]json.RawMessage, comparison StringComparison) error {
	// Validate that the input data matches the input struct
	if err := m.Compare(input); err != nil {
		// Marshal the struct input
		inputModelBytes, jsonErr := json.Marshal(m)
		if jsonErr != nil {
			return fmt.Errorf("failed to marshal input model during comparison error: %s (original error: %s)", jsonErr, err)
		}

		inputBytes, jsonErr := json.Marshal(input)
		if jsonErr != nil {
			return fmt.Errorf("failed to marshal input data during comparison error: %s (original error: %s)", jsonErr, err)
		}

		return fmt.Errorf("input data validation failed: %v, expected struct format: %v, received: %v", err, string(inputModelBytes), string(inputBytes))
	}

	m.ensureCaches() // Ensure caches are initialized

	// Build the path to the attribute if not already cached
	if _, exists := m.attributePaths[comparison.AttributeId]; !exists {
		m.buildAttributePath(comparison.AttributeId)
	}

	// Get the path for this attribute
	path := m.attributePaths[comparison.AttributeId]
	if len(path) == 0 {
		return fmt.Errorf("attribute \"%s\" not found in model", comparison.AttributeId)
	}

	// Get the nested value using the cached path
	val, exists := m.getNestedValue(input, path)
	if !exists {
		return fmt.Errorf("attribute \"%s\" not found in input", comparison.AttributeId)
	}

	valStr, ok := val.(string)
	if !ok {
		return fmt.Errorf("attribute '%s' (%s) is not a string", path, comparison.AttributeId)
	}

	switch comparison.Operator {
	case "equal-to":
		if valStr == comparison.Value {
			return nil
		} else {
			return fmt.Errorf("attribute %s is not equal to %v", comparison.AttributeId, comparison.Value)
		}
	case "contains":
		if strings.Contains(valStr, comparison.Value) {
			return nil
		} else {
			return fmt.Errorf("attribute %s does not contain %v", comparison.AttributeId, comparison.Value)
		}
	case "starts-with":
		if strings.HasPrefix(valStr, comparison.Value) {
			return nil
		} else {
			return fmt.Errorf("attribute %s does not start with %v", comparison.AttributeId, comparison.Value)
		}
	case "ends-with":
		if strings.HasSuffix(valStr, comparison.Value) {
			return nil
		} else {
			return fmt.Errorf("attribute %s does not end with %v", comparison.AttributeId, comparison.Value)
		}
	default:
		return fmt.Errorf("unsupported operator %s", comparison.Operator)
	}
}
