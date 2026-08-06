package certificatesGet

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	gruenthttp "github.com/ongruent/gruent/libs/http"
	config "github.com/ongruent/gruent/services/http-api/config"
	"go.opentelemetry.io/otel/attribute"
)

func execute(request *gruenthttp.Request[Input]) (*Output, error) {
	request.Span.AddEvent("execute")
	request.Span.SetAttributes(
		attribute.String("request.Input.CertificateId", request.Input.CertificateId),
		attribute.String("request.Input.SpaceId", request.Input.SpaceId),
	)

	cert, key, err := config.SecretsClient.RetrieveCert(request.Input.CertificateId)
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "secret_client"),
			attribute.String("error.message", fmt.Sprintf("could not get the certificate: %s", err)),
		)

		errorResponse := gruenthttp.NewErrorResponse([]any{"Could not retrieve the certificate at this time."})
		errorResponseJson, _ := errorResponse.MarshalJson()
		return nil, fiber.NewError(fiber.StatusBadRequest, string(errorResponseJson))
	}

	caCert, err := config.SecretsClient.RetrieveRootCaCert()
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "secret_client"),
			attribute.String("error.message", fmt.Sprintf("could not get the CA certificate: %s", err)),
		)

		errorResponse := gruenthttp.NewErrorResponse([]any{"Could not retrieve the certificate at this time."})
		errorResponseJson, _ := errorResponse.MarshalJson()
		return nil, fiber.NewError(fiber.StatusBadRequest, string(errorResponseJson))
	}

	output := Output{
		Certificate: cert,
		PrivateKey:  key,
		CaCert:      caCert,
	}

	return &output, nil
}
