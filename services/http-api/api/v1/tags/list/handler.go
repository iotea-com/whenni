package tagsList

import (
	"github.com/gofiber/fiber/v2"
	sqldb "github.com/iotea-com/iotea/db/sqlc"
)

type ExpandedTag struct {
	Tag      sqldb.AppTag `json:"tag"`
	Things   []string     `json:"things"`
	Channels []string     `json:"channels"`
	Models   []string     `json:"models"`
}

type Input struct {
	BearerToken string `validate:"required"`
	SpaceId     string `validate:"required"`
	Category    string
}

type Output struct {
	Tags map[string]ExpandedTag `json:"tags"`
}

// @Summary List tags
// @Description Given a space ID, retrieves the names of all tags in the space.
// @Tags tags
// @Accept  json
// @Produce  json
// @Param   spaceId	query	string	true	"Space ID"
// @Param   category	query	string	true	"Category"
// @Success 200
// @Router /tags [get]
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
