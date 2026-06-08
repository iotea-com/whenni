package spacesUpdate

import (
	"github.com/gofiber/fiber/v2"
	sqldb "github.com/iotea-com/iotea/db/sqlc"
)

type SpacePayload struct {
	ID   string `json:"id" validate:"required"`
	Name string `json:"name" validate:"required,min=3,max=30"`
}

type Input struct {
	BearerToken string       `validate:"required"`
	OrgId       string       `validate:"required"`
	Space       SpacePayload `validate:"required"`
}

type Output struct {
	Space *sqldb.AppSpace `json:"space"`
}

// @Summary Update a space
// @Description Given a full space object, update the details of that space.
// @Tags spaces
// @Accept  json
// @Produce  json
// @Param   spaceId	path	string	true	"Space ID"
// @Param   orgId	query	string	true	"Organization ID"
// @Success 200
// @Router /spaces/{spaceId} [put]
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
