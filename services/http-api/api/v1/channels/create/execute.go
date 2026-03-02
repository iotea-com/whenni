package channelsCreate

import (
	"encoding/json"
	"fmt"

	ioteahttp "github.com/iotea-com/iotea/libs/http"
	"github.com/iotea-com/iotea/libs/id"
	"github.com/iotea-com/iotea/libs/legacy/engine/dependencies/things"
	"github.com/iotea-com/iotea/prisma/db"
	"github.com/iotea-com/iotea/services/http-api/services/prisma"
	"github.com/iotea-com/iotea/services/http-api/util"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func execute(request *ioteahttp.Request[Input]) (*Output, error) {
	request.Span.AddEvent("execute")
	request.Span.SetAttributes(
		attribute.String("request.Input.Name", request.Input.Name),
		attribute.String("request.Input.SpaceId", request.Input.SpaceId),
	)

	// Create channel entry
	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Insert channel")
	channelId, err := id.Generator.NewChannelId()
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "id_generation"),
			attribute.String("error.message", fmt.Sprintf("error generating a new Channel ID: %s", err)),
		)
		return nil, err
	}

	if request.Input.Config == nil {
		request.Input.Config = make(map[string]any)
	}
	request.Input.Config["id"] = channelId
	jsonConfig, _ := json.Marshal(request.Input.Config)

	channel, err := prisma.Client.Channel.CreateOne(
		db.Channel.ID.Set(*channelId),
		db.Channel.Name.Set(request.Input.Name),
		db.Channel.Config.Set(jsonConfig),
		db.Channel.CreatedBy.Set(request.GetActorId()),
		db.Channel.UpdatedBy.Set(request.GetActorId()),
		db.Channel.Space.Link(
			db.Space.ID.Equals(request.Input.SpaceId),
		),
	).Exec(dbCtx)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error inserting channel into the database: %s", err)),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	// Create draft internal HTTP server thing
	dbCtx, dbSpan = otel.Tracer("prisma").Start(request.Context, "Insert internal HTTP server thing")

	// Set the thing ID to the same ID than the channel
	thingId := channelId
	now := util.GetCurrentTime()

	_, err = prisma.Client.Thing.CreateOne(
		db.Thing.ID.Set(*thingId),
		db.Thing.Name.Set(request.Input.Name),
		db.Thing.Attributes.Set(db.JSON{'{', '}'}),
		db.Thing.CreatedBy.Set("internal"),
		db.Thing.UpdatedBy.Set("internal"),
		db.Thing.ThingCategory.Set(things.HttpServerThingCategory.String()),
		db.Thing.Space.Link(
			db.Space.ID.Equals(request.Input.SpaceId),
		),
		db.Thing.Internal.Set(true),
		db.Thing.CreatedAt.Set(now),
		db.Thing.UpdatedAt.Set(now),
	).Exec(dbCtx)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error inserting internal HTTP server thing into the database: %s", err)),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	output := Output{
		Channel: channel,
	}

	return &output, nil
}
