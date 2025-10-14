package val

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var secretKeyRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

func IsValidSecretName(fl validator.FieldLevel) bool {
	return secretKeyRegex.MatchString(fl.Field().String())
}
