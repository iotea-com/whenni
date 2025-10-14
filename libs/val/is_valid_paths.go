package val

import (
	"net/url"

	"github.com/go-playground/validator/v10"
)

// Validates that the provided string is a valid URL path
func IsValidPath(fl validator.FieldLevel) bool {
	path := fl.Field().Interface().(string)

	u, err := url.ParseRequestURI(path)
	if err != nil || u.Scheme != "" || u.Host != "" || u.Path != path {
		return false
	}

	return true
}

// Validates that every string in a provided []string is a valid URL path
func IsValidPaths(fl validator.FieldLevel) bool {
	paths := fl.Field().Interface().([]string)

	for _, path := range paths {
		u, err := url.ParseRequestURI(path)
		if err != nil || u.Scheme != "" || u.Host != "" || u.Path != path {
			return false
		}
	}

	return true
}
