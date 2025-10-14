package organizationsMembersRemove

import (
	"github.com/gofiber/fiber/v2"
)

type Input struct {
	BearerToken string `validate:"required"`
	OrgId       string `validate:"required"`
	UserId      string `validate:"required"`
}

type Output struct{}

// @Summary Remove a member from an organization
// @Description Given a valid organization ID and user ID, removes a member from an organization.
// @Tags organizations
// @Accept  json
// @Produce  json
// @Param   orgId	path	string	true	"Organization ID"
// @Param request body organizationsMembersRemove.parse.RequestBody true "Remove user request body"
// @Success 200
// @Router /organizations/{orgId}/members [delete]
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
