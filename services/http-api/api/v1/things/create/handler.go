package thingsCreate

import (
	"github.com/gofiber/fiber/v2"
	sqldb "github.com/ongruent/gruent/db/sqlc"
	"github.com/ongruent/gruent/libs/legacy/engine/dependencies/things"
)

type Input struct {
	BearerToken string               `validate:"required"`
	SpaceId     string               `validate:"required"`
	Name        string               `validate:"required,min=3,max=50"`
	Category    things.ThingCategory `validate:"required"`
	Attributes  map[string]any       `validate:"required"`
}

type Output struct {
	Thing *sqldb.AppThing `json:"thing"`
}

// @Summary Create a thing
// @Description Given a space ID, name, category, and attributes, creates a new thing in a space.
// @Tags things
// @Accept  json
// @Produce  json
// @Param request body thingsCreate.parse.RequestBody true "Body"
// @Param   spaceId	query	string	true	"Space ID"
// @Success 201
// @Router /things [post]
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
