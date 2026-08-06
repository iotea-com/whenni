package tagsRemove

import (
	gruenthttp "github.com/ongruent/gruent/libs/http"
	gruentpermissions "github.com/ongruent/gruent/libs/http/permissions"
	"github.com/ongruent/gruent/services/http-api/config"
	"go.opentelemetry.io/otel/attribute"
)

func contextValidate(request *gruenthttp.Request[Input]) error {
	request.Span.AddEvent("contextValidate")

	authorizeRequestParams := gruenthttp.AuthorizeRequestParams{
		BearerToken: request.Input.BearerToken,
		JwtSecret:   config.VaultConf.JwtSecret,
		ScopeId:     request.Input.SpaceId,
		Namespace:   gruentpermissions.NamespaceSecrets,
		Action:      gruentpermissions.ActionCreate,
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
