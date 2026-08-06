package policiesUpdate

import (
	"github.com/gofiber/fiber/v2"
	sqldb "github.com/ongruent/gruent/db/sqlc"
	mqttPolicies "github.com/ongruent/gruent/libs/http/policies"
)

type Input struct {
	BearerToken   string              `validate:"required"`
	SpaceId       string              `validate:"required"`
	CertificateId string              `validate:"required"`
	Policy        mqttPolicies.Policy `validate:"required"`
	Revoke        bool
}

type Output struct {
	Certificate *sqldb.AppCertificate
}

// @Summary Update an x.509 policy
// @Description Given a certificate ID and a policy, update the policy for a certificate.
// @Tags policies
// @Accept  json
// @Produce  json
// @Param   spaceId	query	string	true	"Space ID"
// @Param   certificateId	path	string	true	"Certificate ID"
// @Param request body policiesUpdate.parse.RequestBody true "Body"
// @Success 200
// @Router /policies/{certificateId} [put]
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
