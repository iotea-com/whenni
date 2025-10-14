package secretsUpdate

import (
	"github.com/gofiber/fiber/v2"
)

type Input struct {
	BearerToken string `validate:"required"`
	Name        string `validate:"required,isValidSecretName,min=3,max=255"`
	Value       string `validate:"required,min=1,max=9999"`
	SpaceId     string `validate:"required"`
}

type Output struct{}

// @Summary Update a secret
// @Description Given a secret name, value, and space ID, updates the secret.
// @Tags secrets
// @Accept  json
// @Produce  json
// @Param request body secretsUpdate.parse.RequestBody true "Body"
// @Param   spaceId	query	string	true	"Space ID"
// @Success 201
// @Router /secrets [post]
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
