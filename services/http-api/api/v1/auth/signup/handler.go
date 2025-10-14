package signup

import (
	"github.com/gofiber/fiber/v2"
)

type Input struct {
	Email       string `validate:"required"`
	Password    string `validate:"required"`
	AppUrl      string `validate:"required"`
	InviteToken *string
}

type Output struct{}

// @Summary Signup
// @Description Creates a new user account.
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   email		body	string	true	"Email"
// @Param   password		body	string	true	"Password"
// @Param   appUrl		body	string	true	"App URL"
// @Param   inviteToken		body	string	false	"Invite Token"
// @Success 200
// @Router /auth/signup [post]
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

	go action(&request.Input, ctx.Context())

	return nil
}
