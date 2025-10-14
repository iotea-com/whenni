package spacesDelete

import (
	"github.com/gofiber/fiber/v2"
)

type Input struct {
	BearerToken string `validate:"required"`
	OrgId       string `validate:"required"`
	SpaceId     string `validate:"required"`
}

type Output struct{}

// @Summary Delete a space
// @Description Given a space ID, deletes the space. All channels must be unpublished before the space can be deleted.
// @Tags spaces
// @Accept  json
// @Produce  json
// @Param   spaceId	path	string	true	"Space ID"
// @Param   orgId	query	string	true	"Organization ID"
// @Success 200
// @Router /spaces/{spaceId} [delete]
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
