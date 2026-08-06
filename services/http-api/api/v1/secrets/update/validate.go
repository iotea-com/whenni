package secretsUpdate

import (
	"github.com/gofiber/fiber/v2"
	gruenthttp "github.com/ongruent/gruent/libs/http"
	"github.com/ongruent/gruent/libs/val"
	"go.opentelemetry.io/otel/attribute"

	validator "github.com/go-playground/validator/v10"
)

func validate(request *gruenthttp.Request[Input]) error {
	request.Span.AddEvent("validate")

	v := validator.New()
	v.RegisterValidation("isValidSecretName", val.IsValidSecretName)

	if err := v.Struct(request.Input); err != nil {
		request.Span.SetAttributes(
			attribute.String("validate.result", "fail"),
			attribute.String("error.type", "input_validation"),
			attribute.String("error.message", err.Error()),
		)
		validationError := gruenthttp.NewGruentValidationError(err.(validator.ValidationErrors))
		response := validationError.MarshalResponse()
		responseJson, err := response.MarshalJson()
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		return fiber.NewError(fiber.StatusBadRequest, string(responseJson))
	}

	request.Span.SetAttributes(attribute.String("validate.result", "pass"))
	return nil
}
