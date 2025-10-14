package secretsDelete

import (
	"github.com/gofiber/fiber/v2"
)

type Input struct {
	BearerToken string `validate:"required"`
	Name        string `validate:"required"`
	SpaceId     string `validate:"required"`
}

type Output struct{}

// @Summary Delete a secret
// @Description Given a space ID and secret name, deletes the secret.
// @Tags secrets
// @Accept  json
// @Produce  json
// @Param   spaceId	query	string	true	"Space ID"
// @Param   secretName	path	string	true	"Secret Name"
// @Success 200
// @Router /secrets/{secretName} [delete]
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
