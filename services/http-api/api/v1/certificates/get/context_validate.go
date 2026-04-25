package certificatesGet

import (
	"go.opentelemetry.io/otel/attribute"

	ioteahttp "github.com/iotea-com/iotea/libs/http"
	ioteapermissions "github.com/iotea-com/iotea/libs/http/permissions"
	"github.com/iotea-com/iotea/services/http-api/config"
)

func contextValidate(request *ioteahttp.Request[Input]) error {
	request.Span.AddEvent("contextValidate")

	authorizeRequestParams := ioteahttp.AuthorizeRequestParams{
		BearerToken: request.Input.BearerToken,
		JwtSecret:   config.VaultConf.JwtSecret,
		ScopeId:     request.Input.SpaceId,
		Namespace:   ioteapermissions.NamespaceCertificates,
		Action:      ioteapermissions.ActionGet,
	}

	err := request.Authorize(authorizeRequestParams)
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "authorize"),
			attribute.String("error.message", err.Error()),
			attribute.String("context_validation.status", "fail"),
		)
		return err
	}

	request.Span.SetAttributes(attribute.String("context_validation.status", "pass"))
	return nil
}
