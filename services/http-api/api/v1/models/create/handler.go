package modelsCreate

import (
	"github.com/gofiber/fiber/v2"
	"github.com/iotea-com/iotea/libs/engine/dependencies/models"
	"github.com/iotea-com/iotea/prisma/db"
)

type Input struct {
	SpaceId     string                      `validate:"required"`
	BearerToken string                      `validate:"required"`
	Name        string                      `validate:"required,min=3,max=100"`
	Attributes  map[string]models.Attribute `validate:"required"`
}

type Output struct {
	Model *db.ModelModel `json:"model"`
}

// @Summary Create a model
// @Description Given a space ID, a name, and a schema, create a new model.
// @Tags models
// @Accept  json
// @Produce  json
// @Param   spaceId	query	string	true	"Space ID"
// @Param request body modelsCreate.parse.RequestBody true "Body"
// @Success 200
// @Router /models [post]
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
