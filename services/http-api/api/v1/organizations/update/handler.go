package organizationsUpdate

import (
	"github.com/gofiber/fiber/v2"
	"github.com/iotea-com/iotea/prisma/db"
)

type Input struct {
	BearerToken  string               `validate:"required"`
	Organization db.OrganizationModel `validate:"required"`
}

type Output struct {
	Organization *db.OrganizationModel `json:"organization"`
}

// @Summary Update an organization
// @Description Given a full organization object, update the details of that organization.
// @Tags organizations
// @Accept  json
// @Produce  json
// @Param request body organizationsUpdate.parse.RequestBody true "Body"
// @Param   orgId	path	string	true	"Organization ID"
// @Success 200
// @Router /organizations/{orgId} [put]
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
