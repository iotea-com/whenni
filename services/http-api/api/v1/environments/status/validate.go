package environmentsStatus

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
)

func validate(request *ioteahttp.Request[Input]) error {
	request.Span.AddEvent("validate")

	v := validator.New()

	if err := v.Struct(request.Input); err != nil {
		validationError := ioteahttp.NewIoteaValidationError(err.(validator.ValidationErrors))
		response := validationError.MarshalResponse()
		responseJson, err := response.MarshalJson()
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		return fiber.NewError(fiber.StatusBadRequest, string(responseJson))
	}

	return nil
}
