package secretsList

import (
	"github.com/gofiber/fiber/v2"
)

type Input struct {
	BearerToken string `validate:"required"`
	SpaceId     string `validate:"required"`
}

type Output struct {
	Secrets []Secret `json:"secrets"`
}

// @Summary List secrets
// @Description Given a space ID, retrieves the names of all secrets in the space.
// @Tags secrets
// @Accept  json
// @Produce  json
// @Param   spaceId	query	string	true	"Space ID"
// @Success 200
// @Router /secrets [get]
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
