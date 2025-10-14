package channelsValidate

import (
	"github.com/gofiber/fiber/v2"
	"github.com/iotea-com/iotea/libs/engine/channels"
)

type Input struct {
	BearerToken string           `validate:"required"`
	SpaceId     string           `validate:"required"`
	Config      channels.Channel `validate:"required"`
}

type Output struct {
	ChannelErrors []string            `json:"channelErrors"`
	NodeErrors    map[string][]string `json:"nodeErrors"`
}

// @Summary Validate a channel configuration
// @Description Given a channel ID and a channel configuration, validate that the configuration is valid.
// @Tags channels
// @Accept  json
// @Produce  json
// @Param   spaceId	query	string	true	"Space ID"
// @Param   channelId	path	string	true	"Channel ID"
// @Success 200
// @Router /channels/{channelId}/validate [patch]
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
