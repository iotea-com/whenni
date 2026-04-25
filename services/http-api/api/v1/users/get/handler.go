package get

import (
	"github.com/gofiber/fiber/v2"
	sqldb "github.com/iotea-com/iotea/db/sqlc"
)

type Input struct {
	UserId string `validate:"required"`
}

type Output struct {
	User *sqldb.AppUser
}

// @Summary Get a user
// @Description Given a user ID, retrieves the user's details.
// @Tags users
// @Accept  json
// @Produce  json
// @Param   userId	path	string	true	"User ID"
// @Success 200
// @Router /users/{userId} [get]
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
