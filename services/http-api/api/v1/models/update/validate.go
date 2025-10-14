package modelsUpdate

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/iotea-com/iotea/libs/engine/dependencies/models"
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

	// Validate model attributes
	var modelAttributes map[string]models.Attribute
	err := json.Unmarshal(request.Input.Model.Attributes, &modelAttributes)
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("validate.result", "fail"),
			attribute.String("error.type", "schema_validation"),
			attribute.String("error.message", err.Error()),
		)
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	model := models.NewModel(modelAttributes)
	err = model.Validate()
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

	request.Input.attributes = modelAttributes

	// TODO: SQL escape JSON before storing it in DB

	request.Span.SetAttributes(attribute.String("validate.result", "pass"))
	return nil
}
