package environmentsSsh

import (
	"github.com/gofiber/fiber/v2"
)

type Input struct {
	SpaceId       string
	BearerToken   string
	EnvironmentID string
}

type Output struct {
	PublicDNS string `json:"publicDNS"`
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
