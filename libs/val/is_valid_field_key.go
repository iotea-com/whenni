package val

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var keyRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*$`)

func IsValidFieldKey(fl validator.FieldLevel) bool {
	match := keyRegex.MatchString(fl.Field().String())
	return match
}
