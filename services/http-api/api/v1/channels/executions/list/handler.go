package list

import (
	"github.com/gofiber/fiber/v2"
)

type Input struct {
	BearerToken    string
	SpaceId        string `validate:"required"`
	ChannelId      string `validate:"required"`
	Page           int    `validate:"gt=0"`
	ResultsPerPage int    `validate:"gt=0,lte=100"`
	StatusFilter   string
}

type Output struct {
	ChannelExecutionsList []ChannelExecutionListItem
	Page                  int
	TotalPages            int
	TotalResults          int
	ResultsPerPage        int
}

// @Summary List channel executions
// @Description Given a channel ID, list the channel executions for that channel.
// @Tags channels, channel-executions
// @Accept  json
// @Produce  json
// @Param   spaceId	query	string	true	"Space ID"
// @Param   channelId	path	string	true	"Channel ID"
// @Param   page	query	int	false	"Page"
// @Param   statusFilter	query	string	false	"Status filter"
// @Param   resultsPerPage	query	int	false	"Results per page"
// @Success 200
// @Router /channels/{channelId}/executions [get]
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
