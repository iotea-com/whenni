package spacesDelete

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
		attribute.String("request.Input.SpaceId", request.Input.SpaceId),
	)

	spaceId := request.Input.SpaceId

	// Get published channels for this space from the database
	noPublishedChannels := false

	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Get channels")
	channels, err := prisma.Client.Channel.FindMany(
		db.Channel.SpaceID.Equals(spaceId),
	).Exec(dbCtx)
	if err != nil {
		if err.Error() == "ErrNotFound" {
			noPublishedChannels = true
		} else {
			dbSpan.SetAttributes(
				attribute.String("error.type", "database"),
				attribute.String("error.message", fmt.Sprintf("error getting channels from the database: %s", err)),
			)
			dbSpan.End()
			return nil, err
		}
	}

	publishedChannels := []string{}
	for _, channel := range channels {
		_, published := channel.PublishedAt()
		if published {
			publishedChannels = append(publishedChannels, channel.Name)
		}
	}

	if len(publishedChannels) <= 0 {
		noPublishedChannels = true
	}

	dbSpan.End()

	// Do not delete space if there are published channels
	if !noPublishedChannels {
		request.Span.SetAttributes(
			attribute.String("error.type", "context_validation"),
			attribute.String("error.message", "delete space failed since there are published channels linked to the space"),
		)
		response := ioteahttp.NewErrorResponse([]any{"unpublish or delete all channels in this space before deleting the space"})
		responseJson, _ := response.MarshalJson()
		return nil, fiber.NewError(fiber.StatusConflict, string(responseJson))
	}

	// Delete models
	dbCtx, dbSpan = otel.Tracer("prisma").Start(request.Context, "Delete models")
	_, err = prisma.Client.Model.FindMany(
		db.Model.SpaceID.Equals(spaceId),
	).Delete().Exec(dbCtx)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error deleting models from the database: %s", err)),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	// Delete things
	dbCtx, dbSpan = otel.Tracer("prisma").Start(request.Context, "Delete things")
	_, err = prisma.Client.Thing.FindMany(
		db.Thing.SpaceID.Equals(spaceId),
	).Delete().Exec(dbCtx)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error deleting things from the database: %s", err)),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	// Delete certificates
	dbCtx, dbSpan = otel.Tracer("prisma").Start(request.Context, "Delete certificates")
	_, err = prisma.Client.Certificate.FindMany(
		db.Certificate.SpaceID.Equals(spaceId),
	).Delete().Exec(dbCtx)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error deleting certificates from the database: %s", err)),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	// Delete api keys
	dbCtx, dbSpan = otel.Tracer("prisma").Start(request.Context, "Delete API keys")
	_, err = prisma.Client.APIKey.FindMany(
		db.APIKey.SpaceID.Equals(spaceId),
	).Delete().Exec(dbCtx)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error deleting API keys from the database: %s", err)),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	// Delete permission sets
	dbCtx, dbSpan = otel.Tracer("prisma").Start(request.Context, "Delete permission sets")
	_, err = prisma.Client.PermissionSet.FindMany(
		db.PermissionSet.SpaceID.Equals(spaceId),
	).Delete().Exec(dbCtx)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error deleting permission sets from the database: %s", err)),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	// Delete channels
	dbCtx, dbSpan = otel.Tracer("prisma").Start(request.Context, "Delete channel")
	_, err = prisma.Client.Channel.FindMany(
		db.Channel.SpaceID.Equals(request.Input.SpaceId),
	).Delete().Exec(dbCtx)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error deleting channels from the database: %s", err)),
		)
		dbSpan.End()

		return nil, err
	}

	dbSpan.End()

	// Delete space
	dbCtx, dbSpan = otel.Tracer("prisma").Start(request.Context, "Delete space")
	_, err = prisma.Client.Space.FindUnique(
		db.Space.ID.Equals(spaceId),
	).Delete().Exec(dbCtx)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error deleting space from the database: %s", err)),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	output := Output{}

	return &output, nil
}
