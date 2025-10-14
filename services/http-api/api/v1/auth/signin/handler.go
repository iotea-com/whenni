package signin

import (
	"github.com/gofiber/fiber/v2"
)

type Input struct {
	Method     string `validate:"required,oneof=credentials magicLink"`
	Email      string `validate:"required"`
	Password   string `validate:"required_if=Method credentials"`
	AppUrl     string `validate:"required_if=Method magicLink"`
	RedirectTo string
}

type Output struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// @Summary Signin
// @Description Given a signin method and relevant credentials, signs the user in.
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   method		body	string	true	"Signin method"
// @Param   email		body	string	true	"Email"
// @Param   password		body	string	false	"Password"
// @Param   appUrl		body	string	true	"App URL"
// @Param   redirectTo		body	string	false	"Redirect to"
// @Success 200
// @Router /auth/signin [post]
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
