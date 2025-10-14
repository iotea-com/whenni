package models

import (
	"encoding/json"
	"fmt"
)

type Condition struct {
	AttributeId string  `json:"attributeId" validate:"required"`
	Operator    string  `json:"operator" validate:"required,oneof=greater-than less-than equal-to less-than-equal-to greater-than-equal-to not-equal-to"`
	Value       float64 `json:"value" validate:"required"`
}

func (m *Model) EvaluateCondition(input map[string]json.RawMessage, condition Condition) error {
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

	// Convert the value to float64 for comparison
	var floatVal float64
	switch v := val.(type) {
	case float64:
		floatVal = v
	case json.Number:
		var err error
		floatVal, err = v.Float64()
		if err != nil {
			return fmt.Errorf("error converting field %s to number: %v", condition.AttributeId, err)
		}
	default:
		return fmt.Errorf("attribute %s is not a number", condition.AttributeId)
	}

	switch condition.Operator {
	case "greater-than":
		if floatVal > condition.Value {
			return nil // Condition met
		} else {
			return fmt.Errorf("attribute %s is not greater than %v", condition.AttributeId, condition.Value)
		}
	case "less-than":
		if floatVal < condition.Value {
			return nil // Condition met
		} else {
			return fmt.Errorf("attribute %s is not less than %v", condition.AttributeId, condition.Value)
		}
	case "equal-to":
		if floatVal == condition.Value {
			return nil // Condition met
		} else {
			return fmt.Errorf("attribute %s is not equal to %v", condition.AttributeId, condition.Value)
		}
	case "less-than-equal-to":
		if floatVal <= condition.Value {
			return nil // Condition met
		} else {
			return fmt.Errorf("attribute %s is not less than or equal to %v", condition.AttributeId, condition.Value)
		}
	case "greater-than-equal-to":
		if floatVal >= condition.Value {
			return nil // Condition met
		} else {
			return fmt.Errorf("attribute %s is not greater than or equal to %v", condition.AttributeId, condition.Value)
		}
	case "not-equal-to":
		if floatVal != condition.Value {
			return nil // Condition met
		} else {
			return fmt.Errorf("attribute %s is equal to %v", condition.AttributeId, condition.Value)
		}
	default:
		return fmt.Errorf("unsupported operator %s", condition.Operator)
	}
}
