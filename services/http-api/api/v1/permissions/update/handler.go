package permissionsUpdate

import (
	"github.com/gofiber/fiber/v2"
	sqldb "github.com/iotea-com/iotea/db/sqlc"
)

type Input struct {
	BearerToken     string `validate:"required"`
	OrgId           string `validate:"required"`
	SpaceId         string
	PermissionSetId string   `validate:"required"`
	Name            string   `validate:"required,min=3,max=50"`
	Permissions     []string `validate:"required"`
}

type Output struct {
	PermissionSet *sqldb.AppPermission
}

// @Summary Update a permission set
// @Description Given a permission set and a list of permissions, update an existing permission set. If a space ID is provided, the provided permissions are expected to be space permissions.
// @Tags permissions
// @Accept  json
// @Produce  json
// @Param   orgId	query	string	true	"Organization ID"
// @Param   spaceId	query	string	false	"Space ID"
// @Param   permissionSetId	path	string	true	"Permission Set ID"
// @Param request body permissionsUpdate.parse.RequestBody true "Body"
// @Success 200
// @Router /permissions/{permissionSetId} [put]
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
