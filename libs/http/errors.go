package ioteahttp

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// A normalized object for input validation errors. Contains a slice of human-readable error
// messages. Can be converted to an API response object.
type IoteaValidationError struct {
	Errors []string
}

// Creates a custom error message mapping based on `validator` validation errors.
func NewIoteaValidationError(validationErrors validator.ValidationErrors) IoteaValidationError {
	ioteaValidationError := IoteaValidationError{
		Errors: []string{},
	}

	// custom error message mapping
	for _, ve := range validationErrors {
		switch ve.Tag() {
		case "required":
			e := fmt.Sprintf("Missing '%s'. This is a required field.", ve.Field())
			ioteaValidationError.Errors = append(ioteaValidationError.Errors, e)
		case "min":
			if ve.Type().Name() == "string" {
				e := fmt.Sprintf("Invalid '%s'. Try again with at least %s characters.", ve.Field(), ve.Param())
				ioteaValidationError.Errors = append(ioteaValidationError.Errors, e)
			}
		case "max":
			e := fmt.Sprintf("Invalid '%s'. Try again with less than %s characters.", ve.Field(), ve.Param())
			ioteaValidationError.Errors = append(ioteaValidationError.Errors, e)
		default:
			e := fmt.Sprintf("Invalid '%s'.", ve.Field())
			ioteaValidationError.Errors = append(ioteaValidationError.Errors, e)
		}
	}

	return ioteaValidationError
}

// Converts the validation error object to an IOTEA API response object.
func (r IoteaValidationError) MarshalResponse() IoteaApiResponse {
	var errors []any
	for _, err := range r.Errors {
		errors = append(errors, err)
	}

	apiResponse := NewErrorResponse(errors)

	return apiResponse
}

// Creates a generic error API response. Sets data to null and errors to a slice of anything
// (ideally human-readable strings).
func NewErrorResponse(errors []any) IoteaApiResponse {
	r := IoteaApiResponse{
		Data:   nil,
		Errors: errors,
	}

	return r
}
