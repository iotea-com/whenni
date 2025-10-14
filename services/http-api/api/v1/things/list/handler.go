package thingsList

import (
	"github.com/gofiber/fiber/v2"
	"github.com/iotea-com/iotea/libs/engine/dependencies/things"
	"github.com/iotea-com/iotea/prisma/db"
)

type Input struct {
	BearerToken    string `validate:"required"`
	SpaceId        string `validate:"required"`
	ThingCategory  *things.ThingCategory
	Page           int `validate:"gt=0"`
	ResultsPerPage int `validate:"gt=0,lte=100"`
	Filter         string
	TagFilter      []string
}

type Output struct {
	Things         []db.ThingModel
	Page           int
	TotalPages     int
	TotalResults   int
	ResultsPerPage int
}

// @Summary List things
// @Description Given a space, list the things in a space. Optionally filter by thing category.
// @Tags things
// @Accept  json
// @Produce  json
// @Param   spaceId	query	string	true	"Space ID"
// @Param   thingCategory	query	string	false	"Thing Category"
// @Param   page	query	int	false	"Page"
// @Param   resultsPerPage	query	int	false	"Results per page"
// @Success 200
// @Router /things [get]
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
