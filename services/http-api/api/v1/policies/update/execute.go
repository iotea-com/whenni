package policiesUpdate

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	"github.com/iotea-com/iotea/services/http-api/services/sqlc"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func execute(request *ioteahttp.Request[Input]) (*Output, error) {
	request.Span.AddEvent("execute")

	// Store in the database
	policyJson, _ := json.Marshal(request.Input.Policy)

	dbCtx, dbSpan := otel.Tracer("sqlc").Start(request.Context, "Update permission set")
	certificate, err := sqlc.Queries.UpdatePolicy(
		dbCtx,
		request.Input.CertificateId,
		request.Input.SpaceId,
		policyJson,
		request.Input.Revoke,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			dbSpan.SetAttributes(
				attribute.String("error.type", "database"),
				attribute.String("error.message", fmt.Sprintf("no certificate found in space with ID %s", request.Input.SpaceId)),
			)
			dbSpan.End()
			return nil, fiber.NewError(fiber.StatusBadRequest)
		}
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error updating permission set in the database: %s", err)),
		)
		dbSpan.End()
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	dbSpan.End()

	output := Output{
		Certificate: &certificate,
	}

	return &output, nil
}
