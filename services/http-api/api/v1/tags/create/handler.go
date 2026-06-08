package tagsCreate

import (
	"github.com/gofiber/fiber/v2"
	sqldb "github.com/iotea-com/iotea/db/sqlc"
)

type Input struct {
	BearerToken string `validate:"required"`
	Name        string `validate:"required,min=1,max=30"`
	SpaceId     string `validate:"required"`
}

type Output struct {
	Tag sqldb.AppTag `json:"tag"`
}

// @Summary Create a tag
// @Description Given a name and a space ID, creates a new tag.
// @Tags tags
// @Accept  json
// @Produce  json
// @Param request body tagsCreate.parse.RequestBody true "Body"
// @Param   spaceId	query	string	true	"Space ID"
// @Success 201
// @Router /tags [post]
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
