package policiesUpdate

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	gruenthttp "github.com/ongruent/gruent/libs/http"
	"github.com/ongruent/gruent/libs/val"
	"go.opentelemetry.io/otel/attribute"
)

func validate(request *gruenthttp.Request[Input]) error {
	request.Span.AddEvent("validate")

	v := validator.New()
	v.RegisterValidation("mqtt_topic", val.MqttTopic)

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

	// Make sure every topic is prefixed with the space ID
	for _, topic := range request.Input.Policy.AllowedSubscriptionTopics {
		if !strings.HasPrefix(topic, request.Input.SpaceId) {
			return fmt.Errorf("subscription topic %s must be prefixed with space ID %s", topic, request.Input.SpaceId)
		}
	}

	for _, topic := range request.Input.Policy.AllowedPublishTopics {
		if !strings.HasPrefix(topic, request.Input.SpaceId) {
			return fmt.Errorf("publish topic %s must be prefixed with space ID %s", topic, request.Input.SpaceId)
		}
	}

	request.Span.SetAttributes(attribute.String("validate.result", "pass"))
	return nil
}
