package val

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

func MqttTopic(fl validator.FieldLevel) bool {
	topic := fl.Field().String()

	// Check if empty
	if len(topic) == 0 {
		return false
	}

	// Check for null characters
	if strings.Contains(topic, "\u0000") {
		return false
	}

	// Split into levels
	levels := strings.Split(topic, "/")

	for _, level := range levels {
		// Empty level (e.g., double slashes)
		if len(level) == 0 {
			return false
		}

		// Single level wildcard validation
		if level == "+" {
			continue
		}

		// Multi level wildcard validation - must be last level
		if level == "#" {
			if level != levels[len(levels)-1] {
				return false
			}
			continue
		}
	}

	return true
}
