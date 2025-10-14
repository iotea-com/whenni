package environmentsCreate

import (
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	ioteapermissions "github.com/iotea-com/iotea/libs/http/permissions"
	"github.com/iotea-com/iotea/services/http-api/config"
	"github.com/iotea-com/iotea/services/http-api/services/prisma"
)

func contextValidate(request *ioteahttp.Request[Input]) error {
	request.Span.AddEvent("contextValidate")

	authorizeRequestParams := ioteahttp.AuthorizeRequestParams{
		BearerToken:  request.Input.BearerToken,
		PrismaClient: prisma.Client,
		JwtSecret:    config.VaultConf.JwtSecret,
		ScopeId:      request.Input.OrgId,
		Namespace:    ioteapermissions.NamespaceEnvironments,
		Action:       ioteapermissions.ActionCreate,
	}

	err := request.Authorize(authorizeRequestParams)
	if err != nil {
		return err
	}

	return nil
}
