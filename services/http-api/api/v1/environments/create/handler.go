package environmentsCreate

import (
	"github.com/gofiber/fiber/v2"
)

type Input struct {
	OrgId       string
	BearerToken string
	Name        string
	Type        string
	Region      string
}

type Output struct {
	EnvironmentID string `json:"environmentId"`
	Message       string `json:"message"`
	PrivateKey    []byte `json:"privateKey"`
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
