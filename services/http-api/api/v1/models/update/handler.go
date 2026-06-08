package modelsUpdate

import (
	"github.com/gofiber/fiber/v2"
	sqldb "github.com/iotea-com/iotea/db/sqlc"
	"github.com/iotea-com/iotea/libs/legacy/engine/dependencies/models"
)

type ModelPayload struct {
	Name       string `json:"name"`
	Attributes []byte `json:"attributes"`
}

type Input struct {
	SpaceId     string        `validate:"required"`
	ModelId     string        `validate:"required"`
	BearerToken string        `validate:"required"`
	Model       ModelPayload  `validate:"required"`

	attributes map[string]models.Attribute
}

type Output struct {
	Model *sqldb.AppModel `json:"model"`
}

// @Summary Update a model
// @Description Given a model ID, update the details of that model.
// @Tags models
// @Accept  json
// @Produce  json
// @Param   spaceId	query	string	true	"Space ID"
// @Param   modelId	path	string	true	"Model ID"
// @Success 200
// @Router /models/{modelId} [put]
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
