package modelsList

import (
	"github.com/gofiber/fiber/v2"
	sqldb "github.com/iotea-com/iotea/db/sqlc"
)

type Input struct {
	SpaceId        string
	BearerToken    string
	Page           int `validate:"gt=0"`
	ResultsPerPage int `validate:"gt=0,lte=100"`
	Filter         string
	TagFilter      []string
}

type Output struct {
	Models         []sqldb.AppModel
	Page           int
	TotalPages     int
	TotalResults   int
	ResultsPerPage int
}

// @Summary List models
// @Description Given a space ID, list all of the models in the space.
// @Tags models
// @Accept  json
// @Produce  json
// @Param   spaceId	query	string	true	"Space ID"
// @Param   page	query	int	false	"Page"
// @Param   resultsPerPage	query	int	false	"Results per page"
// @Success 200
// @Router /models [get]
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
