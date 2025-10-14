package policiesGet

import (
	"fmt"

	ioteahttp "github.com/iotea-com/iotea/libs/http"
	"github.com/iotea-com/iotea/prisma/db"
	"github.com/iotea-com/iotea/services/http-api/services/prisma"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func execute(request *ioteahttp.Request[Input]) (*Output, error) {
	request.Span.AddEvent("execute")
	request.Span.SetAttributes(
		attribute.String("request.Input.CertificateId", request.Input.CertificateId),
		attribute.String("request.Input.SpaceId", request.Input.SpaceId),
	)

	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Get policy")
	certificate, err := prisma.Client.Certificate.FindUnique(
		db.Certificate.ID.Equals(request.Input.CertificateId),
	).Exec(dbCtx)
	if err != nil {
		if err.Error() == "ErrNotFound" {
			dbSpan.SetAttributes(
				attribute.String("error.type", "database"),
				attribute.String("error.message", fmt.Sprintf("no certificate found in space with ID %s", request.Input.SpaceId)),
			)
			dbSpan.End()

			output := Output{
				Policy: nil,
			}

			return &output, nil
		}

		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error listing permission sets in the database: %s", err)),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	output := Output{
		Policy: certificate.Policy,
		Revoke: certificate.Revoke,
	}

	return &output, nil
}
