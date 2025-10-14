package models

import (
	"encoding/json"
	"fmt"
)

/*
Compares the schema against a map of data and returns an error if any fields are missing or types are invalid.

Objects and arrays are recursively validated.
*/
func (m *Model) Compare(data map[string]json.RawMessage) error {
	return compareDataToSchema(m, data)
}

/*
Abstracts the comparison logic from the public method to enable recursion.
*/
func compareDataToSchema(m *Model, data map[string]json.RawMessage) error {
	for id, field := range m.Attributes {
		// Skip if the field is a child of another field
		if field.ParentId != "" {
			continue
		}

		// Check if the field exists in the data; if it's missing and not required, continue to the next field
		rawVal, exists := data[field.Key]
		if !exists && field.Required {
			return fmt.Errorf("expected field '%s' (%s) is missing in input data", field.Key, id)
		}

		if !field.Required && !exists {
			continue
		}

		// Unmarshal the raw value
		var val any
		if err := json.Unmarshal(rawVal, &val); err != nil {
			return fmt.Errorf("error unmarshaling field '%s' (%s): %v", field.Key, id, err)
		}

		// Simple type checking based on the unmarshaled value
		switch field.Type {
		case "string":
			_, isString := val.(string)
			if !isString {
				var v any = val
				if field.Protected {
					v = "[PROTECTED_VALUE]"
				}
				return fmt.Errorf("attribute '%s' (%s) is not a string - got %v", field.Key, id, v)
			}
		case "number":
			var v any = val
			if field.Protected {
				v = "[PROTECTED_VALUE]"
			}
			_, isNumber := val.(float64) // JSON numbers are by default unmarshaled as float64 in Go
			if !isNumber {
				return fmt.Errorf("attribute '%s' (%s) is not a number - got %v", field.Key, id, v)
			}
		case "boolean":
			var v any = val
			if field.Protected {
				v = "[PROTECTED_VALUE]"
			}
			_, isBool := val.(bool)
			if !isBool {
				return fmt.Errorf("attribute '%s' (%s) is not a boolean - got %v", field.Key, id, v)
			}
		case "array":
			isValidSlice := false
			switch field.Subtype {
			case "string":
				var stringSlice []string
				if err := json.Unmarshal(rawVal, &stringSlice); err != nil {
					return fmt.Errorf("attribute '%s' (%s) is not a valid string array: %v", field.Key, id, err)
				}
				isValidSlice = true
			case "number":
				var numberSlice []float64
				if err := json.Unmarshal(rawVal, &numberSlice); err != nil {
					return fmt.Errorf("attribute '%s' (%s) is not a valid number array: %v", field.Key, id, err)
				}
				isValidSlice = true
			case "boolean":
				var boolSlice []bool
				if err := json.Unmarshal(rawVal, &boolSlice); err != nil {
					return fmt.Errorf("attribute '%s' (%s) is not a valid boolean array: %v", field.Key, id, err)
				}
				isValidSlice = true
			case "object":
				var objectSlice []map[string]json.RawMessage
				if err := json.Unmarshal(rawVal, &objectSlice); err != nil {
					return fmt.Errorf("attribute '%s' (%s) is not a valid object array: %v", field.Key, id, err)
				}

				if len(objectSlice) > 0 {
					// Create a subschema from the arrayObjectFields
					subFields := make(map[string]Attribute)
					for subId, subField := range field.ArrayObjectAttributes {
						// Create a copy of the field without parent ID for comparison
						subFields[subId] = Attribute{
							Id:                    subField.Id,
							Key:                   subField.Key,
							Type:                  subField.Type,
							Subtype:               subField.Subtype,
							Required:              subField.Required,
							Protected:             subField.Protected,
							ArrayObjectAttributes: subField.ArrayObjectAttributes,
							// ParentId is intentionally omitted for the comparison
						}
					}
					subSchema := NewModel(subFields)

					// Validate each object in the array
					for i, obj := range objectSlice {
						if err := compareDataToSchema(subSchema, obj); err != nil {
							return fmt.Errorf("validation failed for object %d in array %s (%s): %s", i, field.Key, id, err)
						}
					}
				}
			default:
				// Assume the subtype is a field ID, which will be used to validate the array as an object array
				objectFieldId := field.Subtype
				var objectSlice []map[string]json.RawMessage
				if err := json.Unmarshal(rawVal, &objectSlice); err != nil {
					return fmt.Errorf("attribute '%s' (%s) is not a valid object array: %v", field.Key, id, err)
				}

				// Recursively validate each object in the array
				subFields := make(map[string]Attribute)
				for _, field := range m.Attributes {
					if field.ParentId == objectFieldId {
						subFields[field.Id] = Attribute{
							Id:        field.Id,
							Key:       field.Key,
							Type:      field.Type,
							Subtype:   field.Subtype,
							Required:  field.Required,
							Protected: field.Protected,
							ParentId:  "", // Reset the parent ID to ensure it's not used in the recursive comparison
						}
					}
				}
				subSchema := NewModel(subFields)

				if err := compareDataToSchema(subSchema, objectSlice[0]); err != nil {
					return fmt.Errorf("validation failed for an object in array %s (%s): %s", field.Key, id, err)
				}

				isValidSlice = true
			}

			if !isValidSlice {
				return fmt.Errorf("attribute %s (%s) is not an array of type %s", field.Key, id, field.Subtype)
			}
		case "object":
			var o map[string]json.RawMessage
			if err := json.Unmarshal(rawVal, &o); err != nil {
				return fmt.Errorf("attribute '%s' (%s) is not a valid object: %v", field.Key, id, err)
			}

			// Create a subschema with an updated base case (no parent ID) and recursively validate the object
			subFields := make(map[string]Attribute)
			for _, field := range m.Attributes {
				if field.ParentId == id {
					subFields[field.Id] = Attribute{
						Id:        field.Id,
						Key:       field.Key,
						Type:      field.Type,
						Subtype:   field.Subtype,
						Required:  field.Required,
						Protected: field.Protected,
						ParentId:  "", // Reset the parent ID to ensure it's not used in the recursive comparison
					}
				}
			}
			subSchema := NewModel(subFields)
			if err := compareDataToSchema(subSchema, o); err != nil {
				return fmt.Errorf("validation failed for object %s (%s): %s", field.Key, id, err)
			}
		default:
			// Should never happen because of initial validation rules
			return fmt.Errorf("unsupported field type '%s' in configuration struct", field.Type)
		}
	}

	return nil
}
