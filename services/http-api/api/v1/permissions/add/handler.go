package permissionsAdd

import (
	"github.com/gofiber/fiber/v2"
	"github.com/iotea-com/iotea/prisma/db"
)

type Input struct {
	BearerToken string `validate:"required"`
	OrgId       string `validate:"required"`
	SpaceId     string
	Name        string   `validate:"required,min=3,max=50"`
	Permissions []string `validate:"required"`
}

type Output struct {
	PermissionSet *db.PermissionSetModel
}

// @Summary Create a permission set
// @Description Given a name and a list of permissions, create a new permission set. If a space ID is provided, the provided permissions are expected to be space permissions.
// @Tags permissions
// @Accept  json
// @Produce  json
// @Param   spaceId	query	string	false	"Space ID"
// @Param   orgId	query	string	true	"Organization ID"
// @Param request body permissionsAdd.parse.RequestBody true "Body"
// @Success 201
// @Router /permissions [post]
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
