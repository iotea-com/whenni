package updatePassword

import (
	"github.com/gofiber/fiber/v2"
)

type Input struct {
	BearerToken string `validate:"required"`
	NewPassword string `validate:"required"`
}

type Output struct{}

// @Summary Change password
// @Description Updates the user's password.
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   password		body	string	true	"New password"
// @Success 200
// @Router /auth/change-password [patch]
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
