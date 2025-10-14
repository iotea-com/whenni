package thingsDelete

import (
	"github.com/gofiber/fiber/v2"
)

type Input struct {
	BearerToken string
	SpaceId     string
	ThingId     string
}

type Output struct{}

// @Summary Delete a thing
// @Description Given a thing ID, permanently delete a thing.
// @Tags things
// @Accept  json
// @Produce  json
// @Param   thingId	path	string	true	"Thing ID"
// @Param   spaceId	query	string	true	"Space ID"
// @Success 200
// @Router /things/{thingId} [delete]
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
