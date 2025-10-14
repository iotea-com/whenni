package refresh

import (
	"github.com/gofiber/fiber/v2"
)

type Input struct {
	RefreshToken string `validate:"required"`
}

type Output struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// @Summary Refresh
// @Description Given a refresh token, refreshes the access token.
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   requestBody		body	RequestBody	true	"Refresh token"
// @Success 200
// @Router /auth/refresh [post]
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
