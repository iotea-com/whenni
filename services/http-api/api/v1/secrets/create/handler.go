package secretsCreate

import (
	"github.com/gofiber/fiber/v2"
)

type Input struct {
	BearerToken string `validate:"required"`
	Name        string `validate:"required,min=3,max=255,isValidSecretName"`
	Value       string `validate:"required,min=1,max=9999"`
	SpaceId     string `validate:"required"`
}

type Output struct{}

// @Summary Create a secret
// @Description Given a value and a space ID, creates a new secret.
// @Tags secrets
// @Accept  json
// @Produce  json
// @Param request body secretsCreate.parse.RequestBody true "Body"
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
