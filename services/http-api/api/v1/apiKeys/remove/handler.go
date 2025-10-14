package apiKeysRemove

import (
	"github.com/gofiber/fiber/v2"
)

type Input struct {
	BearerToken string `validate:"required"`
	OrgId       string `validate:"required"`
	SpaceId     string
	ApiKeyId    string `validate:"required"`
}

type Output struct{}

// @Summary Delete an API key
// @Description Given an API key ID, deletes the API key.
// @Tags api-keys
// @Accept  json
// @Produce  json
// @Param   apiKeyId	path	string	true	"API Key ID"
// @Param   orgId		query	string	true	"Organization ID"
// @Param   spaceId		query	string	false	"Space ID"
// @Success 200
// @Router /api-keys/{apiKeyId} [delete]
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
