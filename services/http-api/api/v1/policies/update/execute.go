package policiesUpdate

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	"github.com/iotea-com/iotea/prisma/db"
	"github.com/iotea-com/iotea/services/http-api/services/prisma"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func execute(request *ioteahttp.Request[Input]) (*Output, error) {
	request.Span.AddEvent("execute")

	// Store in the database
	policyJson, _ := json.Marshal(request.Input.Policy)

	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Update permission set")
	certificate, err := prisma.Client.Certificate.FindUnique(
		db.Certificate.ID.Equals(request.Input.CertificateId),
	).Update(
		db.Certificate.Policy.Set(policyJson),
		db.Certificate.Revoke.Set(request.Input.Revoke),
	).Exec(dbCtx)

	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error updating permission set in the database: %s", err)),
		)
		dbSpan.End()
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	dbSpan.End()

	output := Output{
		Certificate: certificate,
	}

	return &output, nil
}
