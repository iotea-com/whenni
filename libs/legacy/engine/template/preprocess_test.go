package template

import (
	"testing"

	"github.com/ongruent/gruent/libs/legacy/engine/dependencies/models"
)

func TestPreprocess(t *testing.T) {
	tests := []struct {
		name        string
		template    Template
		inputData   map[string]any
		expected    string
		notExpected string
		expectedErr string
		expectError bool
	}{
		{
			name: "Doesn't overwrite partials",
			template: Template{
				Value:         "Example uuid(",
				ModelInput:    &models.Model{},
				DefaultValues: map[string]any{},
			},
			inputData:   map[string]any{},
			expected:    "Example uuid(",
			expectError: false,
		},
		{
			name: "Valid template using uuid()",
			template: Template{
				Value:         "Example uuid()",
				ModelInput:    &models.Model{},
				DefaultValues: map[string]any{},
			},
			inputData:   map[string]any{},
			notExpected: "Example uuid()",
			expectError: false,
		},
		{
			name: "Valid template using date()",
			template: Template{
				Value:         "Example date()",
				ModelInput:    &models.Model{},
				DefaultValues: map[string]any{},
			},
			inputData:   map[string]any{},
			notExpected: "Example date()",
			expectError: false,
		},
		{
			name: "Valid template using timestamp()",
			template: Template{
				Value:         "Example timestamp()",
				ModelInput:    &models.Model{},
				DefaultValues: map[string]any{},
			},
			inputData:   map[string]any{},
			notExpected: "Example timestamp()",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.template.Preprocess()

			result := tt.template.Value

			t.Logf("template value after preprocessing: %s", result)

			if tt.expected != "" {
				if result == tt.notExpected {
					t.Errorf("expected result to be: %s, got %s", tt.notExpected, result)
				}
				return
			}

			if tt.notExpected != "" {
				if result == tt.notExpected {
					t.Errorf("expected result to not be: %s", tt.notExpected)
				}
				return
			}

			t.Errorf("invalid expectation - expected and notExpected are both undefined")
		})
	}
}
