package organizationsMembersList

import (
	"github.com/gofiber/fiber/v2"
	sqldb "github.com/iotea-com/iotea/db/sqlc"
)

type Input struct {
	BearerToken    string `validate:"required"`
	OrgId          string `validate:"required"`
	Page           int    `validate:"gt=0"`
	ResultsPerPage int    `validate:"gt=0,lte=100"`
	Filter         string
}

type Output struct {
	OrganizationMembers []sqldb.ListOrganizationMembersRow `json:"organizationMembers"`
	Page                int
	TotalPages          int
	TotalResults        int
	ResultsPerPage      int
}

// @Summary List the members of an organization
// @Description Given a valid organization ID, view the members of the organization.
// @Tags organizations
// @Produce  json
// @Param   orgId	path	string	true	"Organization ID"
// @Param   page	query	int	false	"Page"
// @Param   resultsPerPage	query	int	false	"Results per page"
// @Success 200
// @Router /organizations/{orgId}/members [get]
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
