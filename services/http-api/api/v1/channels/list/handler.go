package channelsList

import (
	"github.com/gofiber/fiber/v2"
	"github.com/iotea-com/iotea/prisma/db"
)

type Input struct {
	BearerToken    string `validate:"required"`
	SpaceId        string `validate:"required"`
	Page           int    `validate:"gt=0"`
	ResultsPerPage int    `validate:"gt=0,lte=100"`
	Filter         string
	TagFilter      []string
}

type Output struct {
	Channels       []db.ChannelModel
	Page           int
	TotalPages     int
	TotalResults   int
	ResultsPerPage int
}

// @Summary List channels
// @Description Given a space ID, list the channels.
// @Tags channels
// @Accept  json
// @Produce  json
// @Param   spaceId	query	string	true	"Space ID"
// @Param   page	query	int	false	"Page"
// @Param   resultsPerPage	query	int	false	"Results per page"
// @Param   filter	query	string	false	"Filter"
// @Param   tagFilter	query	string	false	"Tag Filter"
// @Success 200
// @Router /channels [get]
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
