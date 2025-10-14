package channelsUpdateConfig

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	"github.com/iotea-com/iotea/prisma/db"
	"github.com/iotea-com/iotea/services/http-api/services/prisma"
	"github.com/iotea-com/iotea/services/http-api/util"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func execute(request *ioteahttp.Request[Input]) (*Output, error) {
	request.Span.AddEvent("execute")
	request.Span.SetAttributes(
		attribute.String("request.Input.ChannelId", request.Input.ChannelId),
		attribute.String("request.Input.SpaceId", request.Input.SpaceId),
	)

	now := util.GetCurrentTime()
	jsonConfig, _ := json.Marshal(request.Input.Config)

	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Update channel")
	channel, err := prisma.Client.Channel.FindUnique(
		db.Channel.ID.Equals(request.Input.ChannelId),
	).Update(
		db.Channel.Config.Set(jsonConfig),
		db.Channel.UpdatedAt.Set(now),
	).Exec(dbCtx)
	if err != nil {
		if err.Error() == "ErrNotFound" {
			dbSpan.SetAttributes(
				attribute.String("error.type", "database"),
				attribute.String("error.message", fmt.Sprintf("no channel found with ID %s", request.Input.ChannelId)),
			)
			dbSpan.End()
			return nil, fiber.NewError(fiber.StatusBadRequest)
		}

		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error updating channel in the database: %s", err)),
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
