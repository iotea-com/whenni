package template

import (
	"testing"

	"github.com/ongruent/gruent/libs/legacy/engine/dependencies/models"
)

func TestVerify(t *testing.T) {
	tests := []struct {
		name        string
		template    Template
		expectedErr string
		expectError bool
	}{
		{
			name: "Valid template with correct input types",
			template: Template{
				Value: "{field1} is {field2}",
				ModelInput: &models.Model{
					Attributes: map[string]models.Attribute{
						"field1": {Id: "field1", Key: "field1", Type: "string"},
						"field2": {Id: "field2", Key: "field2", Type: "integer"},
					},
				},
				DefaultValues: map[string]any{
					"field2": 123,
				},
			},
			expectError: false,
		},
		{
			name: "Empty template string",
			template: Template{
				Value: "",
				ModelInput: &models.Model{
					Attributes: map[string]models.Attribute{
						"field1": {Id: "field1", Key: "field1", Type: "string"},
					},
				},
				DefaultValues: map[string]any{"field1": "some value"},
			},
			expectedErr: "template string value cannot be empty",
			expectError: true,
		},
		{
			name: "DefaultValues contains a field not in StructInput",
			template: Template{
				Value: "{field1}",
				ModelInput: &models.Model{
					Attributes: map[string]models.Attribute{
						"field1": {Id: "field1", Key: "field1", Type: "string"},
					},
				},
				DefaultValues: map[string]any{
					"field2": "some value",
				},
			},
			expectedErr: "field 'field2' in DefaultValues does not exist in StructInput",
			expectError: true,
		},
		{
			name: "DefaultValues contains incorrect type (expecting integer, got string)",
			template: Template{
				Value: "{field1} is {field2}",
				ModelInput: &models.Model{
					Attributes: map[string]models.Attribute{
						"field1": {Id: "field1", Key: "field1", Type: "string"},
						"field2": {Id: "field2", Key: "field2", Type: "integer"},
					},
				},
				DefaultValues: map[string]any{
					"field2": "not an integer",
				},
			},
			expectedErr: "field 'field2' expected to be an integer in StructInput, but got value: not an integer",
			expectError: true,
		},
		{
			name: "DefaultValues contains correct integer type",
			template: Template{
				Value: "{field1} is {field2}",
				ModelInput: &models.Model{
					Attributes: map[string]models.Attribute{
						"field1": {Id: "field1", Key: "field1", Type: "string"},
						"field2": {Id: "field2", Key: "field2", Type: "integer"},
					},
				},
				DefaultValues: map[string]any{
					"field2": 123,
				},
			},
			expectError: false,
		},
		{
			name: "Unsupported field type in StructInput",
			template: Template{
				Value: "{field1}",
				ModelInput: &models.Model{
					Attributes: map[string]models.Attribute{
						"field1": {Id: "field1", Key: "field1", Type: "unsupportedType"},
					},
				},
				DefaultValues: map[string]any{
					"field1": "some value",
				},
			},
			expectedErr: "unsupported field type 'unsupportedType' in StructInput for field 'field1'",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.template.Verify()

			if tt.expectError {
				if err == nil {
					t.Errorf("expected an error, but got none")
				} else if err.Error() != tt.expectedErr {
					t.Errorf("expected error: %s, but got: %s", tt.expectedErr, err.Error())
				}
			} else if err != nil {
				t.Errorf("expected no error, but got: %s", err.Error())
			}
		})
	}
}

func TestVerifyStrict(t *testing.T) {
	tests := []struct {
		name        string
		template    Template
		expectedErr string
		expectError bool
	}{
		{
			name: "Valid template with correct input types",
			template: Template{
				Value: "{field1} is {field2}",
				ModelInput: &models.Model{
					Attributes: map[string]models.Attribute{
						"field1": {Id: "field1", Key: "field1", Type: "string"},
						"field2": {Id: "field2", Key: "field2", Type: "integer"},
					},
				},
				DefaultValues: map[string]any{
					"field2": 123,
				},
			},
			expectError: false,
		},
		{
			name: "Empty template string",
			template: Template{
				Value: "",
				ModelInput: &models.Model{
					Attributes: map[string]models.Attribute{
						"field1": {Id: "field1", Key: "field1", Type: "string"},
					},
				},
				DefaultValues: map[string]any{"field1": "some value"},
			},
			expectedErr: "template string value cannot be empty",
			expectError: true,
		},
		{
			name: "Empty StructInput",
			template: Template{
				Value:         "{field1}",
				ModelInput:    &models.Model{},
				DefaultValues: map[string]any{"field1": "some value"},
			},
			expectedErr: "template input struct cannot be empty",
			expectError: true,
		},
		{
			name: "Missing field in StructInput",
			template: Template{
				Value: "{field1} is {field2}",
				ModelInput: &models.Model{
					Attributes: map[string]models.Attribute{
						"field1": {Id: "field1", Key: "field1", Type: "string"},
					},
				},
				DefaultValues: map[string]any{},
			},
			expectedErr: "not all fields in template were found in struct, missing field in input struct: [field2]",
			expectError: true,
		},
		{
			name: "DefaultValues contains a field not in StructInput",
			template: Template{
				Value: "{field1}",
				ModelInput: &models.Model{
					Attributes: map[string]models.Attribute{
						"field1": {Id: "field1", Key: "field1", Type: "string"},
					},
				},
				DefaultValues: map[string]any{
					"field2": "some value",
				},
			},
			expectedErr: "field 'field2' in DefaultValues does not exist in StructInput",
			expectError: true,
		},
		{
			name: "DefaultValues contains incorrect type (expecting integer, got string)",
			template: Template{
				Value: "{field1} is {field2}",
				ModelInput: &models.Model{
					Attributes: map[string]models.Attribute{
						"field1": {Id: "field1", Key: "field1", Type: "string"},
						"field2": {Id: "field2", Key: "field2", Type: "integer"},
					},
				},
				DefaultValues: map[string]any{
					"field2": "not an integer",
				},
			},
			expectedErr: "field 'field2' expected to be an integer in StructInput, but got value: not an integer",
			expectError: true,
		},
		{
			name: "DefaultValues contains correct integer type",
			template: Template{
				Value: "{field1} is {field2}",
				ModelInput: &models.Model{
					Attributes: map[string]models.Attribute{
						"field1": {Id: "field1", Key: "field1", Type: "string"},
						"field2": {Id: "field2", Key: "field2", Type: "integer"},
					},
				},
				DefaultValues: map[string]any{
					"field2": 123,
				},
			},
			expectError: false,
		},
		{
			name: "Unsupported field type in StructInput",
			template: Template{
				Value: "{field1}",
				ModelInput: &models.Model{
					Attributes: map[string]models.Attribute{
						"field1": {Id: "field1", Key: "field1", Type: "unsupportedType"},
					},
				},
				DefaultValues: map[string]any{
					"field1": "some value",
				},
			},
			expectedErr: "unsupported field type 'unsupportedType' in StructInput for field 'field1'",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.template.VerifyStrict()

			if tt.expectError {
				if err == nil {
					t.Errorf("expected an error, but got none")
				} else if err.Error() != tt.expectedErr {
					t.Errorf("expected error: %s, but got: %s", tt.expectedErr, err.Error())
				}
			} else if err != nil {
				t.Errorf("expected no error, but got: %s", err.Error())
			}
		})
	}
}

func TestFill(t *testing.T) {
	tests := []struct {
		name        string
		template    Template
		inputData   map[string]any
		expected    string
		expectedErr string
		expectError bool
	}{
		{
			name: "Valid template with correct input types",
			template: Template{
				Value: "{field1} is {field2}",
				ModelInput: &models.Model{
					Attributes: map[string]models.Attribute{
						"field1": {Id: "field1", Key: "field1", Type: "string"},
						"field2": {Id: "field2", Key: "field2", Type: "integer"},
					},
				},
				DefaultValues: map[string]any{},
			},
			inputData:   map[string]any{"field1": "John", "field2": 123.0}, // Note: 123 is passed as float64
			expected:    "John is 123",
			expectError: false,
		},
		{
			name: "Use default value for field not in input data",
			template: Template{
				Value: "{field1} is {field2}",
				ModelInput: &models.Model{
					Attributes: map[string]models.Attribute{
						"field1": {Id: "field1", Key: "field1", Type: "string"},
						"field2": {Id: "field2", Key: "field2", Type: "integer"},
					},
				},
				DefaultValues: map[string]any{"field2": 456},
			},
			inputData:   map[string]any{"field1": "Jane"},
			expected:    "Jane is 456",
			expectError: false,
		},
		{
			name: "Replace '{-}' with inputData",
			template: Template{
				Value:         "{-}",
				ModelInput:    &models.Model{},
				DefaultValues: map[string]any{},
			},
			inputData:   map[string]any{"field1": "Jane"},
			expected:    `{"field1":"Jane"}`,
			expectError: false,
		},
		{
			name: "Incorrect type in input data",
			template: Template{
				Value: "{field1} is {field2}",
				ModelInput: &models.Model{
					Attributes: map[string]models.Attribute{
						"field1": {Id: "field1", Key: "field1", Type: "string"},
						"field2": {Id: "field2", Key: "field2", Type: "integer"},
					},
				},
				DefaultValues: map[string]any{},
			},
			inputData:   map[string]any{"field1": "John", "field2": "not an integer"},
			expectedErr: "field 'field2' expected to be an integer, but got value: not an integer (type mismatch)",
			expectError: true,
		},
		{
			name: "Unsupported field type",
			template: Template{
				Value: "{field1}",
				ModelInput: &models.Model{
					Attributes: map[string]models.Attribute{
						"field1": {Id: "field1", Key: "field1", Type: "unsupportedType"},
					},
				},
				DefaultValues: map[string]any{},
			},
			inputData:   map[string]any{"field1": "some value"},
			expectedErr: "unsupported field type 'unsupportedType' for field 'field1'",
			expectError: true,
		},
		{
			name: "Boolean type validation",
			template: Template{
				Value: "{field1} is {field2}",
				ModelInput: &models.Model{
					Attributes: map[string]models.Attribute{
						"field1": {Id: "field1", Key: "field1", Type: "string"},
						"field2": {Id: "field2", Key: "field2", Type: "boolean"},
					},
				},
				DefaultValues: map[string]any{},
			},
			inputData:   map[string]any{"field1": "John", "field2": true},
			expected:    "John is true",
			expectError: false,
		},
		{
			name: "Boolean type mismatch",
			template: Template{
				Value: "{field1} is {field2}",
				ModelInput: &models.Model{
					Attributes: map[string]models.Attribute{
						"field1": {Id: "field1", Key: "field1", Type: "string"},
						"field2": {Id: "field2", Key: "field2", Type: "boolean"},
					},
				},
				DefaultValues: map[string]any{},
			},
			inputData:   map[string]any{"field1": "John", "field2": "not a boolean"},
			expectedErr: "field 'field2' expected to be a boolean, but got value: not a boolean (type mismatch)",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.template.Fill(tt.inputData)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, but got none")
				} else if err.Error() != tt.expectedErr {
					t.Errorf("expected error: %s, but got: %s", tt.expectedErr, err.Error())
				}
			} else if err != nil {
				t.Errorf("expected no error, but got: %s", err.Error())
			} else if result == nil || *result != tt.expected {
				t.Errorf("expected result: %s, but got: %s", tt.expected, *result)
			}
		})
	}
}

func TestFillStrict(t *testing.T) {
	tests := []struct {
		name        string
		template    Template
		inputData   map[string]any
		expected    string
		expectedErr string
		expectError bool
	}{
		{
			name: "Valid template with correct input types",
			template: Template{
				Value: "{field1} is {field2}",
				ModelInput: &models.Model{
					Attributes: map[string]models.Attribute{
						"field1": {Id: "field1", Key: "field1", Type: "string"},
						"field2": {Id: "field2", Key: "field2", Type: "integer"},
					},
				},
				DefaultValues: map[string]any{},
			},
			inputData:   map[string]any{"field1": "John", "field2": 123.0}, // Note: 123 is passed as float64
			expected:    "John is 123",
			expectError: false,
		},
		{
			name: "Use default value for field not in input data",
			template: Template{
				Value: "{field1} is {field2}",
				ModelInput: &models.Model{
					Attributes: map[string]models.Attribute{
						"field1": {Id: "field1", Key: "field1", Type: "string"},
						"field2": {Id: "field2", Key: "field2", Type: "integer"},
					},
				},
				DefaultValues: map[string]any{"field2": 456},
			},
			inputData:   map[string]any{"field1": "Jane"},
			expected:    "Jane is 456",
			expectError: false,
		},
		{
			name: "Replace '{-}' with inputData",
			template: Template{
				Value: "{-}",
				ModelInput: &models.Model{
					Attributes: map[string]models.Attribute{
						"-": {Id: "-", Key: "-", Type: "object"},
					},
				},
				DefaultValues: map[string]any{},
			},
			inputData:   map[string]any{"field1": "Jane"},
			expected:    `{"field1":"Jane"}`,
			expectError: false,
		},
		{
			name: "Incorrect type in input data",
			template: Template{
				Value: "{field1} is {field2}",
				ModelInput: &models.Model{
					Attributes: map[string]models.Attribute{
						"field1": {Id: "field1", Key: "field1", Type: "string"},
						"field2": {Id: "field2", Key: "field2", Type: "integer"},
					},
				},
				DefaultValues: map[string]any{},
			},
			inputData:   map[string]any{"field1": "John", "field2": "not an integer"},
			expectedErr: "field 'field2' expected to be an integer, but got value: not an integer (type mismatch)",
			expectError: true,
		},
		{
			name: "Verify fails due to missing field in StructInput",
			template: Template{
				Value: "{field1} is {field3}",
				ModelInput: &models.Model{
					Attributes: map[string]models.Attribute{
						"field1": {Id: "field1", Key: "field1", Type: "string"},
					},
				},
				DefaultValues: map[string]any{},
			},
			inputData:   map[string]any{"field1": "John"},
			expectedErr: "failed to verify template: not all fields in template were found in struct, missing field in input struct: [field3]",
			expectError: true,
		},
		{
			name: "Unsupported field type",
			template: Template{
				Value: "{field1}",
				ModelInput: &models.Model{
					Attributes: map[string]models.Attribute{
						"field1": {Id: "field1", Key: "field1", Type: "unsupportedType"},
					},
				},
				DefaultValues: map[string]any{},
			},
			inputData:   map[string]any{"field1": "some value"},
			expectedErr: "unsupported field type 'unsupportedType' for field 'field1'",
			expectError: true,
		},
		{
			name: "Boolean type validation",
			template: Template{
				Value: "{field1} is {field2}",
				ModelInput: &models.Model{
					Attributes: map[string]models.Attribute{
						"field1": {Id: "field1", Key: "field1", Type: "string"},
						"field2": {Id: "field2", Key: "field2", Type: "boolean"},
					},
				},
				DefaultValues: map[string]any{},
			},
			inputData:   map[string]any{"field1": "John", "field2": true},
			expected:    "John is true",
			expectError: false,
		},
		{
			name: "Boolean type mismatch",
			template: Template{
				Value: "{field1} is {field2}",
				ModelInput: &models.Model{
					Attributes: map[string]models.Attribute{
						"field1": {Id: "field1", Key: "field1", Type: "string"},
						"field2": {Id: "field2", Key: "field2", Type: "boolean"},
					},
				},
				DefaultValues: map[string]any{},
			},
			inputData:   map[string]any{"field1": "John", "field2": "not a boolean"},
			expectedErr: "field 'field2' expected to be a boolean, but got value: not a boolean (type mismatch)",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.template.FillStrict(tt.inputData)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, but got none")
				} else if err.Error() != tt.expectedErr {
					t.Errorf("expected error: %s, but got: %s", tt.expectedErr, err.Error())
				}
			} else if err != nil {
				t.Errorf("expected no error, but got: %s", err.Error())
			} else if result == nil || *result != tt.expected {
				t.Errorf("expected result: %s, but got: %s", tt.expected, *result)
			}
		})
	}
}
