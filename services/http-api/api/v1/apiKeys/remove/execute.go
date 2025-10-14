package apiKeysRemove

import (
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
	request.Span.SetAttributes(
		attribute.String("request.Input.ApiKeyId", request.Input.ApiKeyId),
		attribute.String("request.Input.OrgId", request.Input.OrgId),
		attribute.String("request.Input.SpaceId", request.Input.SpaceId),
	)

	// Delete the API key from the database
	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Delete API key")
	_, err := prisma.Client.APIKey.FindUnique(
		db.APIKey.ID.Equals(request.Input.ApiKeyId),
	).Delete().Exec(dbCtx)

	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error deleting API key from the database: %s", err)),
		)
		dbSpan.End()
		return nil, fiber.NewError(fiber.StatusInternalServerError)
	}

	dbSpan.End()

	return nil, nil
}
