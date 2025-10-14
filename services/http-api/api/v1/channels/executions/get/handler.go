package get

import (
	"github.com/gofiber/fiber/v2"
)

type Input struct {
	BearerToken        string
	SpaceId            string `validate:"required"`
	ChannelExecutionId string `validate:"required"`
}

type Output struct {
	ChannelExecution *ChannelExecution
}

// @Summary Get a channel execution
// @Description Given a channel execution ID, retrieve the details of that channel execution.
// @Tags channels, channel-executions
// @Accept  json
// @Produce  json
// @Param   spaceId	query	string	true	"Space ID"
// @Param   channelId	path	string	true	"Channel ID"
// @Param   channelExecutionId	path	string	true	"Channel Execution ID"
// @Success 200
// @Router /channels/{channelId}/executions/{channelExecutionId} [get]
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
