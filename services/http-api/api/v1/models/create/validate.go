package modelsCreate

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	"github.com/iotea-com/iotea/libs/legacy/engine/dependencies/models"
	"go.opentelemetry.io/otel/attribute"
)

func validate(request *ioteahttp.Request[Input]) error {
	request.Span.AddEvent("validate")

	// Validate request body
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

	// Validate model attributes
	model := models.NewModel(request.Input.Attributes)
	err := model.Validate()
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("validate.result", "fail"),
			attribute.String("error.type", "schema_validation"),
			attribute.String("error.message", err.Error()),
		)
		validationError := ioteahttp.NewErrorResponse([]any{err.Error()})
		responseJson, err := validationError.MarshalJson()
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		return fiber.NewError(fiber.StatusBadRequest, string(responseJson))
	}

	request.Input.Attributes = model.Attributes

	// TODO: SQL escape JSON before storing it in DB

	request.Span.SetAttributes(attribute.String("validate.result", "pass"))
	return nil
}
