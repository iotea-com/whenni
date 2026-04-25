package organizationsGet

import (
	"github.com/gofiber/fiber/v2"
	sqldb "github.com/iotea-com/iotea/db/sqlc"
)

type Input struct {
	BearerToken string `validate:"required"`
	OrgId       string `validate:"required"`
}

type Output struct {
	Organization *sqldb.AppOrganization `json:"organization"`
}

// @Summary Get an organization
// @Description Given an organization ID, retrieves the details of a specific organization.
// @Tags organizations
// @Accept  json
// @Produce  json
// @Param   orgId	path	string	true	"Organization ID"
// @Success 200
// @Router /organizations/{orgId} [get]
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
