package channelsUnpublish

import (
	"github.com/gofiber/fiber/v2"
)

type Input struct {
	BearerToken string `validate:"required"`
	SpaceId     string `validate:"required"`
	ChannelId   string `validate:"required"`
}

type Output struct{}

// @Summary Unpublish a channel
// @Description Given a channel ID, unpublish a channel. Unpublishing happens asynchronously.
// @Tags channels
// @Accept  json
// @Produce  json
// @Param   spaceId	query	string	true	"Space ID"
// @Param   channelId	path	string	true	"Channel ID"
// @Success 200
// @Router /channels/{channelId}/unpublish [patch]
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
