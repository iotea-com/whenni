package thingsHealthcheck

import (
	"github.com/gofiber/fiber/v2"
)

type Input struct {
	BearerToken   string         `validate:"required"`
	SpaceId       string         `validate:"required"`
	ThingCategory string         `validate:"required"`
	Attributes    map[string]any `validate:"required"`
}

type Output struct{}

// @Summary Healthcheck thing attributes
// @Description Given a thing category and attributes, healthcheck that thing.
// @Tags things
// @Accept  json
// @Produce  json
// @Param   spaceId	query	string	true	"Space ID"
// @Param   thingCategory	query	string	true	"Thing Category"
// @Param   attributes	body	Input	true	"Attributes"
// @Success 200
// @Router /things/healthcheck [patch]
// @Security ApiKeyAuth
func Handler(ctx *fiber.Ctx) error {
	request, err := parse(ctx)
	if err != nil {
		return err
	}

	err = validate(request)
	if err != nil {
		return err
	}

	err = contextValidate(request)
	if err != nil {
		return err
	}

	output, err := execute(request)
	if err != nil {
		return err
	}

	respond(request, output)

	go action(output)

	return nil
}
