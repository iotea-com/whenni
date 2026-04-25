package organizationsUpdate

import (
	"github.com/gofiber/fiber/v2"
	sqldb "github.com/iotea-com/iotea/db/sqlc"
)

type OrganizationPayload struct {
	ID   string `json:"id" validate:"required"`
	Name string `json:"name" validate:"required,min=3,max=30"`
}

type Input struct {
	BearerToken  string              `validate:"required"`
	Organization OrganizationPayload `validate:"required"`
}

type Output struct {
	Organization *sqldb.AppOrganization `json:"organization"`
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
