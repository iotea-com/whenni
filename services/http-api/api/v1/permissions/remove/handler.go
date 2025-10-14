package permissionsRemove

import (
	"github.com/gofiber/fiber/v2"
)

type Input struct {
	BearerToken     string `validate:"required"`
	OrgId           string `validate:"required"`
	PermissionSetId string `validate:"required"`
}

type Output struct{}

// @Summary Remove a permission set
// @Description Given a permission set ID, remove a permission set from a space.
// @Tags permissions
// @Accept  json
// @Produce  json
// @Param   orgId	query	string	true	"Org ID"
// @Param request body permissionsRemove.parse.RequestBody true "Body"
// @Success 200
// @Router /permissions [delete]
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
