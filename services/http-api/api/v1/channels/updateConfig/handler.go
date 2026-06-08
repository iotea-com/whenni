package channelsUpdateConfig

import (
	"github.com/gofiber/fiber/v2"
	sqldb "github.com/iotea-com/iotea/db/sqlc"
)

type Input struct {
	BearerToken string         `validate:"required"`
	SpaceId     string         `validate:"required"`
	ChannelId   string         `validate:"required"`
	Config      map[string]any `validate:"required"`
}

type Output struct {
	Channel *sqldb.AppChannel `json:"channel"`
}

// @Summary Update a channel configuration
// @Description Given a channel ID and a configuration, update the channel's configuration. If the channel is currently published, it will need to be unpublished and published again for these changes to take effect.
// @Tags channels
// @Accept  json
// @Produce  json
// @Param   spaceId	query	string	true	"Space ID"
// @Param   channelId	path	string	true	"Channel ID"
// @Param request body map[string]any true "Body"
// @Success 200
// @Router /channels/{channelId}/config [patch]
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
