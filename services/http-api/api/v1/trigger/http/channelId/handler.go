package channelId

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type Input struct {
	ChannelId         string `validate:"required"`
	Method            string `validate:"required"`
	Body              map[string]any
	BinaryData        []byte `json:"binaryData,omitempty"`
	BinaryContentType string `json:"binaryContentType,omitempty"`
}

type Output struct {
	response *http.Response
}

// @Summary Trigger an HTTP source node
// @Description Given a channel ID, triggers an HTTP source node if available.
// @Tags trigger
// @Accept  json
// @Produce  json
// @Param request body map[string]any false "Body"
// @Param   channelId	path	string	true	"Channel ID"
// @Success 200 {object} channelId.Output
// @Router /trigger/http/{channelId} [post]
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
