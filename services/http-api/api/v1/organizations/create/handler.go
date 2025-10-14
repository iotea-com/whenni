package organizationsCreate

import (
	"github.com/gofiber/fiber/v2"
	"github.com/iotea-com/iotea/prisma/db"
)

type Input struct {
	BearerToken string `validate:"required"`
	Name        string `validate:"required,min=3,max=30"`
}

type Output struct {
	Organization *db.OrganizationModel `json:"organization"`
}

// @Summary Create an organization
// @Description Given a name, creates a new organization. Happy connecting!
// @Tags organizations
// @Accept  json
// @Produce  json
// @Param request body organizationsCreate.parse.RequestBody true "Body"
// @Success 200
// @Router /organizations [post]
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
