package models

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/ongruent/gruent/libs/val"
)

/*
Validates the schema according to the following rules:

- Field keys must be alphanumeric and start with a letter

- Field keys must be unique to their level of nesting

- Objects must have at least one field

- Arrays must specify a subtype

- If array subtype is not primitive, it must reference a valid field ID

- Circular references are not allowed
*/
func (m *Model) Validate() error {
	// Validate main fields
	return validateFields(m, "root schema")
}

// Helper function to validate a set of fields
func validateFields(m *Model, parentContext string) error {
	// Create a cache for key uniqueness checks
	keyLevelMap := make(map[string]map[string]bool) // parentId -> set of keys

	for _, field := range m.Attributes {
		// Initial validation
		v := validator.New()
		v.RegisterValidation("isValidFieldKey", val.IsValidFieldKey)
		if err := v.Struct(field); err != nil {
			return fmt.Errorf("invalid field '%s' (%s) in %s: %s", field.Key, field.Id, parentContext, err)
		}

		// Check key uniqueness within same level
		if keyLevelMap[field.ParentId] == nil {
			keyLevelMap[field.ParentId] = make(map[string]bool)
		}
		if keyLevelMap[field.ParentId][field.Key] {
			return fmt.Errorf("duplicate key '%s' (%s) found in the same object '%s'", field.Key, field.Id, field.ParentId)
		}
		keyLevelMap[field.ParentId][field.Key] = true

		// Validate nested arrays and objects
		switch field.Type {
		case "string", "number", "boolean":
			continue
		case "array":
			// Validate array subtype
			switch field.Subtype {
			case "string", "number", "boolean":
				continue
			case "object":
				// Validate array object fields
				subSchema := NewModel(field.ArrayObjectAttributes)

				if field.ArrayObjectAttributes != nil {
					if err := validateFields(subSchema, fmt.Sprintf("array object '%s'", field.Key)); err != nil {
						return err
					}
				}
			default:
				// If not primitive, it must reference a valid field ID
				if _, exists := m.Attributes[field.Subtype]; !exists {
					return fmt.Errorf("field '%s' (%s): invalid array subtype reference '%s'", field.Key, field.Id, field.Subtype)
				}
			}
		case "object":
			if err := m.validateObject(field); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *Model) validateObject(field Attribute) error {
	// Count children in main fields
	children := []string{}
	for _, f := range m.Attributes {
		if f.ParentId == field.Id {
			children = append(children, f.Id)
		}
	}

	// Also count children in array object fields if they exist
	if field.ArrayObjectAttributes != nil {
		for _, f := range field.ArrayObjectAttributes {
			if f.ParentId == field.Id {
				children = append(children, f.Id)
			}
		}
	}

	// Check that the object has at least one child field
	if len(children) <= 0 {
		return fmt.Errorf("field '%s' (%s): objects must have at least one child field", field.Key, field.Id)
	}

	// Check for circular references in subfields
	if err := m.detectCircularRef(field.Id, make(map[string]bool)); err != nil {
		return fmt.Errorf("field '%s' (%s): %s", field.Key, field.Id, err)
	}

	return nil
}

func (m *Model) detectCircularRef(fieldId string, visited map[string]bool) error {
	if visited[fieldId] {
		return fmt.Errorf("circular reference detected")
	}

	visited[fieldId] = true
	field := m.Attributes[fieldId]

	// Check children in main fields if this is an object
	if field.Type == "object" {
		for _, otherField := range m.Attributes {
			if otherField.ParentId == fieldId {
				if err := m.detectCircularRef(otherField.Id, visited); err != nil {
					return err
				}
			}
		}
	}

	// Check children in array object fields if they exist
	if field.Type == "array" && field.Subtype == "object" && field.ArrayObjectAttributes != nil {
		for _, otherField := range field.ArrayObjectAttributes {
			if otherField.ParentId == fieldId {
				if err := m.detectCircularRef(otherField.Id, visited); err != nil {
					return err
				}
			}
		}
	}

	// Check array subtype if it references another field
	if field.Type == "array" && field.Subtype != "" && field.Subtype != "object" {
		if _, exists := m.Attributes[field.Subtype]; exists {
			if err := m.detectCircularRef(field.Subtype, visited); err != nil {
				return err
			}
		}
	}

	delete(visited, fieldId) // Remove from visited when backtracking
	return nil
}
