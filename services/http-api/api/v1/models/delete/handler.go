package modelsDelete

import (
	"github.com/gofiber/fiber/v2"
)

type Input struct {
	SpaceId     string
	ModelId     string
	BearerToken string
}

type Output struct{}

// @Summary Delete a model
// @Description Given a model ID, delete it.
// @Tags models
// @Accept  json
// @Produce  json
// @Param   spaceId	query	string	true	"Space ID"
// @Param   modelId	path	string	true	"Model ID"
// @Success 200
// @Router /models/{modelId} [delete]
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
