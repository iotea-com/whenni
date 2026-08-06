package apiKeysList

import (
	"github.com/gofiber/fiber/v2"
	sqldb "github.com/ongruent/gruent/db/sqlc"
)

type Input struct {
	BearerToken    string `validate:"required"`
	OrgId          string `validate:"required"`
	SpaceId        string
	Page           int `validate:"gt=0"`
	ResultsPerPage int `validate:"gt=0,lte=100"`
}

type Output struct {
	ApiKeys        []sqldb.AppApiKey `json:"apiKeys"`
	Page           int
	TotalPages     int
	TotalResults   int
	ResultsPerPage int
}

// @Summary List API keys
// @Description Given an organization ID and a space ID (if intending to list API keys in a space), lists all of the API keys belonging to the organization or space.
// @Tags api-keys
// @Accept  json
// @Produce  json
// @Param   orgId		query	string	true	"Organization ID"
// @Param   spaceId		query	string	false	"Space ID"
// @Param   page		query	int		false	"Page"
// @Param   resultsPerPage	query	int		false	"Results Per Page"
// @Success 200
// @Router /api-keys [get]
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
