package apiKeysAdd

import (
	"github.com/gofiber/fiber/v2"
	sqldb "github.com/iotea-com/iotea/db/sqlc"
)

type Input struct {
	BearerToken     string `validate:"required"`
	OrgId           string `validate:"required"`
	SpaceId         string
	PermissionSetId string `validate:"required"`
	Name            string `validate:"required,min=3,max=50"`
}

type Output struct {
	ApiKey *sqldb.AppApiKey
}

// @Summary Create an API key
// @Description Given an API key name, a permission set, an organization ID, and a space ID (if intending to create an API key in a space), creates a new API key for programmatic use.
// @Tags api-keys
// @Accept  json
// @Produce  json
// @Param   orgId		query	string	true	"Organization ID"
// @Param   spaceId		query	string	false	"Space ID"
// @Param request body apiKeysAdd.parse.RequestBody true "Body"
// @Success 200
// @Router /api-keys [post]
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
