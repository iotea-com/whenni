package magicLinkVerify

import (
	"github.com/gofiber/fiber/v2"
)

type Input struct {
	VerificationToken string `validate:"required"`
}

type Output struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// @Summary Verify a magic link
// @Description Given a magic link token, verifies the magic link.
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   token		query	string	true	"Magic link token"
// @Success 200
// @Router /auth/magic-link/verify [get]
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
