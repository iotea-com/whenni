package environmentsDelete

import (
	"github.com/gofiber/fiber/v2"
)

type Input struct {
	SpaceId       string
	BearerToken   string
	Region        string
	EnvironmentID string
}

type Output struct {
	Message string `json:"message"`
}

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
