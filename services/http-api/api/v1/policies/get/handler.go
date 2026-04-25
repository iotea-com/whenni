package policiesGet

import (
	"github.com/gofiber/fiber/v2"
	sqldb "github.com/iotea-com/iotea/db/sqlc"
)

type Input struct {
	BearerToken   string `validate:"required"`
	SpaceId       string `validate:"required"`
	CertificateId string `validate:"required"`
}

type Output struct {
	Certificate *sqldb.AppCertificate `json:"certificate"`
}

// @Summary Get an x.509 policy
// @Description Given a certificate ID, retrieves the attached policy for that certificate.
// @Tags policies
// @Accept  json
// @Produce  json
// @Param   spaceId	query	string	true	"Space ID"
// @Param   certificateId	path	string	true	"Certificate ID"
// @Success 200
// @Router /policies/{certificateId} [get]
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
