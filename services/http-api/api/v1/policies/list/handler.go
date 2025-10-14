package policiesList

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
	Certificates   []db.CertificateModel `json:"certificates"`
	Page           int
	TotalPages     int
	TotalResults   int
	ResultsPerPage int
}

// @Summary List x.509 policies
// @Description Given a space ID, list all of the policies in the space.
// @Tags policies
// @Accept  json
// @Produce  json
// @Param   spaceId	query	string	true	"Space ID"
// @Param   page	query	int	false	"Page"
// @Param   resultsPerPage	query	int	false	"Results per page"
// @Success 200
// @Router /policies [get]
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
