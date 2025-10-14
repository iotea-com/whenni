package search

import (
	"github.com/gofiber/fiber/v2"
	"github.com/iotea-com/iotea/prisma/db"
)

type Input struct {
	Query string `validate:"required"`
	OrgId string
}

type Output struct {
	Users []db.UserModel
}

// @Summary Search for users
// @Description Given a search query, searches for users. If an organization ID is provided, it will exclude users in that organization.
// @Tags users
// @Accept  json
// @Produce  json
// @Param   orgId	query	string	false	"Organization ID"
// @Param   q	query	string	true	"Search query"
// @Success 200
// @Router /users/search [post]
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
