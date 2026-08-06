package spacesCreate

import (
	"github.com/gofiber/fiber/v2"
	sqldb "github.com/ongruent/gruent/db/sqlc"
)

type Input struct {
	BearerToken string `validate:"required"`
	Name        string `validate:"required,min=3,max=30"`
	OrgId       string `validate:"required"`
}

type Output struct {
	Space *sqldb.AppSpace `json:"space"`
}

// @Summary Create a space
// @Description Given a name and an organization ID, creates a new space.
// @Tags spaces
// @Accept  json
// @Produce  json
// @Param request body spacesCreate.parse.RequestBody true "Body"
// @Param   orgId	query	string	true	"Organization ID"
// @Success 201
// @Router /spaces [post]
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
