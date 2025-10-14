package thingsUpdate

import (
	"github.com/gofiber/fiber/v2"
	"github.com/iotea-com/iotea/prisma/db"
)

type Input struct {
	BearerToken string         `validate:"required"`
	SpaceId     string         `validate:"required"`
	ThingId     string         `validate:"required"`
	Name        string         `validate:"required"`
	Attributes  map[string]any `validate:"required"`
}

type Output struct {
	Thing *db.ThingModel `json:"thing"`
}

// @Summary Update a thing
// @Description Given a thing ID, name, and attributes, update the details of that thing.
// @Tags things
// @Accept  json
// @Produce  json
// @Param   spaceId query	string	true	"Space ID"
// @Param   thingId	path	string	true	"Thing ID"
// @Param request body thingsUpdate.parse.RequestBody true "Body"
// @Success 200
// @Router /things/{thingId} [put]
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
