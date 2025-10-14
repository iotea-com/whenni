package certificatesGet

import (
	"github.com/gofiber/fiber/v2"
)

type Input struct {
	BearerToken   string `validate:"required"`
	SpaceId       string `validate:"required"`
	CertificateId string `validate:"required"`
}

type Output struct {
	Certificate *string
	PrivateKey  *string
	CaCert      *string
}

// @Summary Get a certificate bundle
// @Description Given a valid space ID and certificate ID, retrieves the certificate bundle.
// @Tags certificates
// @Accept  json
// @Produce  application/zip
// @Param   spaceId	query	string	true	"Space ID"
// @Param   certificateId	path	string	true	"Certificate ID"
// @Success 200
// @Router /certificates/{certificateId} [get]
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
