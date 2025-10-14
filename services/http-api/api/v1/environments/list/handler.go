package environmentsList

import (
	"github.com/gofiber/fiber/v2"
	"github.com/iotea-com/iotea/prisma/db"
)

type Input struct {
	BearerToken    string `validate:"required"`
	SpaceId        string `validate:"required"`
	Page           int    `validate:"gt=0"`
	ResultsPerPage int    `validate:"gt=0,lte=100"`
}

type Output struct {
	EnvironmentsList []db.DevEnvironmentModel `json:"environmentsList"`
	Page             int
	TotalPages       int
	TotalResults     int
	ResultsPerPage   int
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
