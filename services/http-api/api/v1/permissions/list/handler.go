package permissionsList

import (
	"github.com/gofiber/fiber/v2"
	sqldb "github.com/iotea-com/iotea/db/sqlc"
)

type Input struct {
	BearerToken    string `validate:"required"`
	OrgId          string `validate:"required_without=SpaceId"`
	SpaceId        string `validate:"required_without=OrgId"`
	Page           int    `validate:"gt=0"`
	ResultsPerPage int    `validate:"gt=0,lte=100"`
}

type Output struct {
	PermissionSets []sqldb.AppPermission `json:"permissionSets"`
	Page           int                   `json:"page"`
	TotalPages     int                   `json:"totalPages"`
	TotalResults   int                   `json:"totalResults"`
	ResultsPerPage int                   `json:"resultsPerPage"`
}

// @Summary List permission sets
// @Description Given a scope ID, list the permission sets in the scope.
// @Tags permissions
// @Accept  json
// @Produce  json
// @Param   orgId	query	string	false	"Organization ID"
// @Param   spaceId	query	string	false	"Space ID"
// @Param   page	query	int	false	"Page"
// @Param   resultsPerPage	query	int	false	"Results per page"
// @Success 201
// @Router /permissions [get]
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
