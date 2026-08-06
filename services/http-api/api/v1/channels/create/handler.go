package channelsCreate

import (
	"github.com/gofiber/fiber/v2"
	sqldb "github.com/ongruent/gruent/db/sqlc"
)

type Input struct {
	BearerToken string         `validate:"required"`
	SpaceId     string         `validate:"required"`
	Name        string         `validate:"required,min=3,max=50"`
	Config      map[string]any `validate:"required"`
}

type Output struct {
	Channel *sqldb.AppChannel `json:"channel"`
}

// @Summary Create a channel
// @Description Given a channel name and a configuration, create a new channel.
// @Tags channels
// @Accept  json
// @Produce  json
// @Param   spaceId	query	string	true	"Space ID"
// @Param request body channelsCreate.parse.RequestBody true "Body"
// @Success 201
// @Router /channels [post]
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
