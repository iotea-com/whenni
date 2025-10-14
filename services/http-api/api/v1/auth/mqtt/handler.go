package mqtt

import (
	"github.com/gofiber/fiber/v2"
)

type Input struct {
	CertString string
	Topic      string
	Action     string
}

type Output struct {
	Result string // "allow" | "deny"
}

func Handler(ctx *fiber.Ctx) error {
	request, err := parse(ctx)
	if request == nil {
		return nil
	}
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
	if output == nil {
		return nil
	}
	if err != nil {
		return err
	}

	respond(request, output)

	go action(output)

	return nil
}
