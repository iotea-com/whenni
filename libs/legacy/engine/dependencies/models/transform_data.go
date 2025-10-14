package models

import (
	"encoding/json"
	"fmt"
)

/*
TransformData transforms the input data to the outputSchema using the provided mapping.
The input data is validated against the input schema. If any fields are missing, a default
value is set based on the output field type.

The mapping is a map of output field IDs to input field IDs.

**This implementation supports up to 20 levels of nested objects.** Our platform does not currently support deeper nesting.
*/
func (m *Model) TransformData(outputSchema *Model, input map[string]json.RawMessage, mapping map[string]string) (map[string]any, map[string]any, error) {
	m.ensureCaches()
	outputSchema.ensureCaches()

	// Validate that the input data matches the input struct
	if err := m.Compare(input); err != nil {
		// Marshal the struct input
		inputSchemaBytes, jsonErr := json.Marshal(m)
		if jsonErr != nil {
			return nil, nil, fmt.Errorf("failed to marshal input schema during comparison error: %s (original error: %s)", jsonErr, err)
		}

		inputBytes, jsonErr := json.Marshal(input)
		if jsonErr != nil {
			return nil, nil, fmt.Errorf("failed to marshal input data during comparison error: %s (original error: %s)", jsonErr, err)
		}

		return nil, nil, fmt.Errorf("input data validation failed: %v, expected struct format: %v, received: %v", err, string(inputSchemaBytes), string(inputBytes))
	}

	// Use cached paths or build them
	for _, inFieldId := range mapping {
		if _, exists := m.attributePaths[inFieldId]; !exists {
			m.buildAttributePath(inFieldId)
		}
	}

	// Pre-unmarshal and cache parent objects
	for _, path := range m.attributePaths {
		if len(path) > 1 {
			parentKey := path[0]
			if _, exists := m.parentObjects[parentKey]; !exists {
				if rawValue, ok := input[parentKey]; ok {
					var nestedMap map[string]json.RawMessage
					if err := json.Unmarshal(rawValue, &nestedMap); err != nil {
						continue
					}
					m.parentObjects[parentKey] = nestedMap
				}
			}
		}
	}

	// Create a flat map of field values first
	flatOutputMap := make(map[string]any, len(mapping))
	for outFieldId, inFieldId := range mapping {
		path := m.attributePaths[inFieldId]
		currentValue, exists := m.getNestedValue(input, path)
		if !exists {
			// Field is missing, determine type from struct and set default
			fieldType := outputSchema.Attributes[outFieldId].Type
			var defaultValue any

			switch fieldType {
			case "number":
				defaultValue = 0
			case "boolean":
				defaultValue = false
			case "string":
				defaultValue = ""
			case "array":
				defaultValue = []any{}
			case "object":
				defaultValue = map[string]any{}
			default:
				defaultValue = nil
			}

			flatOutputMap[outFieldId] = defaultValue
		} else {
			flatOutputMap[outFieldId] = currentValue
		}
	}

	// Build the nested output structure
	outputData, maskedOutputData := outputSchema.buildNestedOutput(flatOutputMap)
	return outputData, maskedOutputData, nil
}

// Updated getNestedValue with cache
func (m *Model) getNestedValue(data map[string]json.RawMessage, path []string) (any, bool) {
	if len(path) == 0 {
		return nil, false
	}

	currentKey := path[0]
	rawValue, exists := data[currentKey]
	if !exists {
		return nil, false
	}

	if len(path) == 1 {
		var result any
		if err := json.Unmarshal(rawValue, &result); err != nil {
			return nil, false
		}
		return result, true
	}

	// For nested paths, unmarshal the current level into a map
	var nestedData map[string]json.RawMessage
	if err := json.Unmarshal(rawValue, &nestedData); err != nil {
		return nil, false
	}

	// Recursively process the remaining path
	return m.getNestedValue(nestedData, path[1:])
}

// Optimize buildNestedOutput with pre-allocated maps and fewer type assertions
func (m *Model) buildNestedOutput(flatMap map[string]any) (map[string]any, map[string]any) {
	// Initialize result maps
	unprotectedResult := make(map[string]any)
	protectedResult := make(map[string]any)

	// Helper function to get or create the full path to a nested object
	getOrCreatePath := func(current map[string]any, path []string) map[string]any {
		for _, key := range path {
			next, exists := current[key]
			if !exists {
				newMap := make(map[string]any)
				current[key] = newMap
				current = newMap
			} else if nextMap, ok := next.(map[string]any); ok {
				current = nextMap
			} else {
				// If it exists but isn't a map, replace it with a map
				newMap := make(map[string]any)
				current[key] = newMap
				current = newMap
			}
		}
		return current
	}

	// Helper function to build the complete path for a field
	buildFieldPath := func(fieldId string) []string {
		currentId := fieldId

		var path []string
		for currentId != "" {
			if field, exists := m.Attributes[currentId]; exists {
				path = append([]string{field.Key}, path...)
				currentId = field.ParentId
			} else {
				break
			}
		}
		return path
	}

	// First create all object structures
	for fieldId, field := range m.Attributes {
		if field.Type == "object" {
			path := buildFieldPath(fieldId)
			if len(path) > 0 {
				if field.ParentId == "" {
					// Root level object
					unprotectedResult[field.Key] = make(map[string]any)
					if field.Protected {
						protectedResult[field.Key] = "[PROTECTED_VALUE]"
					} else {
						protectedResult[field.Key] = make(map[string]any)
					}
				} else {
					// Nested object
					parentPath := path[:len(path)-1]
					if len(parentPath) > 0 {
						unprotectedParent := getOrCreatePath(unprotectedResult, parentPath)
						unprotectedParent[field.Key] = make(map[string]any)

						protectedParent := getOrCreatePath(protectedResult, parentPath)
						if field.Protected {
							protectedParent[field.Key] = "[PROTECTED_VALUE]"
						} else {
							protectedParent[field.Key] = make(map[string]any)
						}
					}
				}
			}
		}
	}

	// Then populate all field values
	for fieldId, field := range m.Attributes {
		if field.Type != "object" {
			value, exists := flatMap[fieldId]
			if !exists {
				continue
			}

			path := buildFieldPath(fieldId)
			if len(path) > 0 {
				if field.ParentId == "" {
					// Root level field
					unprotectedResult[field.Key] = value
					if field.Protected {
						protectedResult[field.Key] = "[PROTECTED_VALUE]"
					} else {
						protectedResult[field.Key] = value
					}
				} else {
					// Nested field
					parentPath := path[:len(path)-1]
					if len(parentPath) > 0 {
						// Always set unprotected value
						unprotectedParent := getOrCreatePath(unprotectedResult, parentPath)
						unprotectedParent[field.Key] = value

						// For protected output, check if any parent is protected
						isParentProtected := false
						currentId := field.ParentId
						for currentId != "" {
							if parent, exists := m.Attributes[currentId]; exists {
								if parent.Protected {
									isParentProtected = true
									break
								}
								currentId = parent.ParentId
							} else {
								break
							}
						}

						protectedParent := getOrCreatePath(protectedResult, parentPath)
						if field.Protected || isParentProtected {
							protectedParent[field.Key] = "[PROTECTED_VALUE]"
						} else {
							protectedParent[field.Key] = value
						}
					}
				}
			}
		}
	}

	return unprotectedResult, protectedResult
}

// buildAttributePath traverses the schema to build an array of keys leading to the target field
func (m *Model) buildAttributePath(fieldId string) []string {
	m.cacheBuildMutex.RLock()
	if path, exists := m.attributePaths[fieldId]; exists {
		m.cacheBuildMutex.RUnlock()
		return path
	}
	m.cacheBuildMutex.RUnlock()

	path := []string{}
	field, exists := m.Attributes[fieldId]
	if !exists {
		return path
	}

	path = append(path, field.Key)
	currentField := field
	for currentField.ParentId != "" {
		if parent, exists := m.Attributes[currentField.ParentId]; exists {
			path = append([]string{parent.Key}, path...)
			currentField = parent
		} else {
			break
		}
	}

	m.cacheBuildMutex.Lock()
	m.attributePaths[fieldId] = path
	m.cacheBuildMutex.Unlock()
	return path
}

// Ensure that empty caches are initialized
func (m *Model) ensureCaches() {
	if m.attributePaths == nil {
		m.attributePaths = make(map[string][]string)
	}
	if m.parentObjects == nil {
		m.parentObjects = make(map[string]map[string]json.RawMessage)
	}
	if m.nestedMaps == nil {
		m.nestedMaps = make(map[string]map[string]any)
	}
}
