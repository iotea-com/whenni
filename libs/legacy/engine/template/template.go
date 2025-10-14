package template

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	"github.com/iotea-com/iotea/libs/engine/dependencies/models"
)

// Template represents a structure used to dynamically replace placeholders
// in a template string. It allows verification of input data types and provides
// default values for any missing attributes.
//
// Fields:
//   - Value: The template string containing placeholders in the format {attribute_name}.
//   - ModelInput: A map where keys are the attribute names expected in the template, and
//     values represent the expected data type ("string", "integer", "float", "boolean").
//   - DefaultValues: A map where keys are the attribute names, and values are the default values
//     to use if no corresponding value is provided in the input data.
type Template struct {
	Value         string
	ModelInput    *models.Model
	DefaultValues map[string]any
}

// VerifyStrict ensures that the template is valid and all the necessary attributes are present
// in both the ModelInput and DefaultValues maps. It checks that:
// 1. The template string is not empty.
// 2. ModelInput is provided with expected attribute names and types.
// 3. All placeholders in the template string correspond to a attribute in ModelInput.
// 4. Any default values match the expected type from ModelInput.
//
// Returns an error if:
// - The template string or ModelInput is empty.
// - Any placeholders in the template are not present in ModelInput.
// - A attribute in DefaultValues is not listed in ModelInput or its type does not match.
func (t *Template) VerifyStrict() error {
	// Check that all the attributes are present
	if len(t.Value) < 1 {
		return fmt.Errorf("template string value cannot be empty")
	}

	if len(t.ModelInput.Attributes) < 1 {
		return fmt.Errorf("template input model cannot be empty")
	}

	// Regular expression to capture {var_name} patterns
	re := regexp.MustCompile(`\{(\-|\w+)\}`)
	matches := re.FindAllStringSubmatch(t.Value, -1)

	// Extract all matched values inside {} and append them to templatedValues
	var templatedValues []string
	for _, match := range matches {
		if len(match) > 1 {
			templatedValues = append(templatedValues, match[1])
		}
	}

	// Verify list against input struct
	valuesFoundCount := 0
	var missingFields []string
	for _, templatedValue := range templatedValues {
		if _, exists := t.ModelInput.Attributes[templatedValue]; exists {
			valuesFoundCount++
		} else {
			missingFields = append(missingFields, templatedValue)
		}
	}

	// If any attributes in the template are missing in StructInput, return an error
	if len(templatedValues) != valuesFoundCount {
		return fmt.Errorf("not all attributes in template were found in struct, missing attribute in input struct: %v", missingFields)
	}

	// Verify that the attributes in DefaultValues exist in StructInput and their types match
	for attributeKey, defaultValue := range t.DefaultValues {
		// Check if the attribute exists in StructInput
		var attribute *models.Attribute
		for _, f := range t.ModelInput.Attributes {
			if f.Key == attributeKey {
				attribute = &f
				break
			}
		}
		if attribute == nil {
			return fmt.Errorf("attribute '%s' in DefaultValues does not exist in StructInput", attributeKey)
		}

		expectedType := attribute.Type

		// Perform type checking based on the expected type
		switch expectedType {
		case "string":
			if _, isString := defaultValue.(string); !isString {
				return fmt.Errorf("attribute '%s' expected to be a string in StructInput, but got value: %v", attributeKey, defaultValue)
			}
		case "float":
			if _, isFloat64 := defaultValue.(float64); !isFloat64 {
				return fmt.Errorf("attribute '%s' expected to be a float in StructInput, but got value: %v", attributeKey, defaultValue)
			}
		case "integer":
			// DefaultValues are stored as string, so we need to convert them to float64 before checking if they can be an integer
			if num, err := convertToFloat64(defaultValue); err != nil || num != float64(int(num)) {
				return fmt.Errorf("attribute '%s' expected to be an integer in StructInput, but got value: %v", attributeKey, defaultValue)
			}
		case "boolean":
			if _, isBool := defaultValue.(bool); !isBool {
				return fmt.Errorf("attribute '%s' expected to be a boolean in StructInput, but got value: %v", attributeKey, defaultValue)
			}
		default:
			return fmt.Errorf("unsupported attribute type '%s' in StructInput for attribute '%s'", expectedType, attributeKey)
		}
	}

	return nil
}

// Verify ensures that the template is valid. It checks that:
// 1. The template string is not empty.
// 2. Any default values match the expected type from StructInput.
//
// Returns an error if:
// - The template string is empty.
// - A attribute in DefaultValues is not listed in StructInput or its type does not match.
func (t *Template) Verify() error {
	// Check that all the attributes are present
	if len(t.Value) < 1 {
		return fmt.Errorf("template string value cannot be empty")
	}

	// Verify that the attributes in DefaultValues exist in StructInput and their types match
	for attributeKey, defaultValue := range t.DefaultValues {
		// Check if the attribute exists in StructInput
		var attribute *models.Attribute
		for _, f := range t.ModelInput.Attributes {
			if f.Key == attributeKey {
				attribute = &f
				break
			}
		}
		if attribute == nil {
			return fmt.Errorf("attribute '%s' in DefaultValues does not exist in StructInput", attributeKey)
		}

		expectedType := attribute.Type

		// Perform type checking based on the expected type
		switch expectedType {
		case "string":
			if _, isString := defaultValue.(string); !isString {
				return fmt.Errorf("attribute '%s' expected to be a string in StructInput, but got value: %v", attributeKey, defaultValue)
			}
		case "float":
			if _, isFloat64 := defaultValue.(float64); !isFloat64 {
				return fmt.Errorf("attribute '%s' expected to be a float in StructInput, but got value: %v", attributeKey, defaultValue)
			}
		case "integer":
			// DefaultValues are stored as string, so we need to convert them to float64 before checking if they can be an integer
			if num, err := convertToFloat64(defaultValue); err != nil || num != float64(int(num)) {
				return fmt.Errorf("attribute '%s' expected to be an integer in StructInput, but got value: %v", attributeKey, defaultValue)
			}
		case "boolean":
			if _, isBool := defaultValue.(bool); !isBool {
				return fmt.Errorf("attribute '%s' expected to be a boolean in StructInput, but got value: %v", attributeKey, defaultValue)
			}
		default:
			return fmt.Errorf("unsupported attribute type '%s' in StructInput for attribute '%s'", expectedType, attributeKey)
		}
	}

	return nil
}

func (t *Template) fill(inputData map[string]any, strict bool) (*string, error) {
	// Always run Verify first
	if strict {
		if err := t.VerifyStrict(); err != nil {
			return nil, fmt.Errorf("failed to verify template: %s", err)
		}
	} else {
		if err := t.Verify(); err != nil {
			return nil, fmt.Errorf("failed to verify template: %s", err)
		}
	}

	// Regular expression to capture {var_name} patterns, including the global key of {-} for
	// original input data
	re := regexp.MustCompile(`\{(\-|\w+)\}`)

	// To store an error in case of a mismatch
	var validationError error

	// Replace each occurrence of {var_name} in the template with the corresponding value from
	// inputData
	filledTemplate := re.ReplaceAllStringFunc(t.Value, func(placeholder string) string {
		attributeKey := re.FindStringSubmatch(placeholder)[1]

		// Use "-" as a global key for the original inputData
		if attributeKey == "-" {
			marshalledInputData, err := json.Marshal(inputData)
			if err != nil {
				validationError = fmt.Errorf("recognized global '-', but could not marshal input data into string: %s", err)
				return placeholder
			}

			return string(marshalledInputData)
		}

		// Get the expected type from StructInput
		// TODO: this should be recursive to find attributes in nested objects
		attribute, exists := t.ModelInput.Attributes[attributeKey]
		if !exists {
			for _, f := range t.ModelInput.Attributes {
				if f.Key == attributeKey {
					attribute = f
					exists = true
					break
				}
			}
		}
		expectedType := attribute.Type

		// Look up the value from inputData
		value, exists := inputData[attributeKey]
		if exists {
			// Perform type checking based on the expected type
			switch expectedType {
			case "string":
				if _, isString := value.(string); !isString {
					validationError = fmt.Errorf("attribute '%s' expected to be a string, but got value: %v (type mismatch)", attributeKey, value)
					return placeholder
				}
			case "float":
				if _, isFloat64 := value.(float64); !isFloat64 {
					validationError = fmt.Errorf("attribute '%s' expected to be a float, but got value: %v (type mismatch)", attributeKey, value)
					return placeholder
				}
			case "integer":
				// Handle float64 as integer if it represents an integer value
				if num, isFloat64 := value.(float64); isFloat64 {
					if num != float64(int(num)) {
						validationError = fmt.Errorf("attribute '%s' expected to be an integer, but got value: %v (type mismatch)", attributeKey, value)
						return placeholder
					}
				} else {
					validationError = fmt.Errorf("attribute '%s' expected to be an integer, but got value: %v (type mismatch)", attributeKey, value)
					return placeholder
				}
			case "boolean":
				if _, isBool := value.(bool); !isBool {
					validationError = fmt.Errorf("attribute '%s' expected to be a boolean, but got value: %v (type mismatch)", attributeKey, value)
					return placeholder
				}
			default:
				validationError = fmt.Errorf("unsupported attribute type '%s' for attribute '%s'", expectedType, attributeKey)
				return placeholder
			}

			// If no error, use fmt.Sprintf with %v to format any type
			return fmt.Sprintf("%v", value)
		}

		// If not found in inputData, look up in DefaultValues
		if defaultValue, exists := t.DefaultValues[attributeKey]; exists {
			return fmt.Sprintf("%v", defaultValue)
		}

		if strict {
			// If neither found, return an error indicating the missing attribute and why
			validationError = fmt.Errorf("attribute '%s' is missing in both inputData and DefaultValues", attributeKey)
		}
		return placeholder
	})

	// If a validation error occurred during the process, return it
	if validationError != nil {
		return nil, validationError
	}

	// Return the filled template as a pointer to the string
	return &filledTemplate, nil
}

// Fill will
// - Run Verify() on the template to ensure that both StructInput and DefaultValues are consistent.
// - Replace all {attribute_name} patterns in the template with their corresponding values from inputData.
// - Fallback to DefaultValues if the attribute is missing from inputData.
// - Return the filled template string or an error.
func (t *Template) Fill(inputData map[string]any) (*string, error) {
	filledTemplate, err := t.fill(inputData, false)
	return filledTemplate, err
}

// Fill will
// - Run VerifyStrict() on the template to ensure that both StructInput and DefaultValues are consistent.
// - Replace all {attribute_name} patterns in the template with their corresponding values from inputData.
// - Fallback to DefaultValues if the attribute is missing from inputData.
// - Return an error if a attribute is missing from both inputData and DefaultValues, providing a clear message.
// - Return the filled template string or an error if any mismatch or inconsistency is found.
func (t *Template) FillStrict(inputData map[string]any) (*string, error) {
	filledTemplate, err := t.fill(inputData, true)
	return filledTemplate, err
}

// Helper function to convert a value to float64
func convertToFloat64(value any) (float64, error) {
	switch v := value.(type) {
	case string:
		return strconv.ParseFloat(v, 64)
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	default:
		return 0, fmt.Errorf("unable to convert value '%v' to float64", value)
	}
}
