package mqtt

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	"go.opentelemetry.io/otel/attribute"
)

func validate(request *ioteahttp.Request[Input]) error {
	request.Span.AddEvent("validate")

	v := validator.New()

	if err := v.Struct(request.Input); err != nil {
		request.Span.SetAttributes(
			attribute.String("validate.result", "fail"),
			attribute.String("error.type", "input_validation"),
			attribute.String("error.message", err.Error()),
		)
		validationError := ioteahttp.NewIoteaValidationError(err.(validator.ValidationErrors))
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
