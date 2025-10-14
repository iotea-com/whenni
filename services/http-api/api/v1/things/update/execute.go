package thingsUpdate

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	ioteachannels "github.com/iotea-com/iotea/libs/engine/channels"
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
		attribute.String("request.Input.Name", request.Input.Name),
		attribute.String("request.Input.ThingId", request.Input.ThingId),
		attribute.String("request.Input.SpaceId", request.Input.SpaceId),
		attribute.String("request.Input.Attributes", fmt.Sprintf("%#v", request.Input.Attributes)),
	)

	// Check if thing is used in any channel
	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Check if thing is used in channel")
	channels, err := prisma.Client.Channel.FindMany(
		db.Channel.SpaceID.Equals(request.Input.SpaceId),
	).Select(
		db.Channel.PublishedAt.Field(),
		db.Channel.Config.Field(),
	).Exec(dbCtx)
	if err != nil && err.Error() != "ErrNotFound" {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error checking if thing is used in a channel: %s", err)),
		)
		dbSpan.End()

		return nil, err
	}

	dbSpan.End()

	channelsWithThing := []db.ChannelModel{}
	for _, channel := range channels {
		// Check if channel is published
		_, isPublished := channel.PublishedAt()
		if !isPublished {
			continue
		}

		// Unmarshal channel config
		var channelConfig ioteachannels.Channel
		err = json.Unmarshal(channel.Config, &channelConfig)
		if err != nil {
			request.Span.SetAttributes(
				attribute.String("error.type", "parse_json"),
				attribute.String("error.message", fmt.Sprintf("could not read channel config: %s", err.Error())),
			)

			errorResponse := ioteahttp.NewErrorResponse([]any{
				fmt.Sprintf("could not read channel config: %s", err.Error()),
			})
			errorResponseJson, _ := errorResponse.MarshalJson()
			return nil, fiber.NewError(fiber.StatusBadRequest, string(errorResponseJson))
		}

		// Check if the thing is used in the channel config
		for _, node := range channelConfig.Nodes {
			for _, thingDependency := range node.Metadata.Dependencies.Things {

				if thingDependency.ThingId == request.Input.ThingId {
					channelsWithThing = append(channelsWithThing, channel)
				}
			}
		}
	}

	if len(channelsWithThing) > 0 {
		request.Span.SetAttributes(
			attribute.String("error.type", "context_validation"),
			attribute.String("error.message", fmt.Sprintf("thing is used in %d channel(s)", len(channelsWithThing))),
		)

		errorResponse := ioteahttp.NewErrorResponse([]any{
			fmt.Sprintf("thing is used in %d channel(s)", len(channelsWithThing)),
		})
		errorResponseJson, _ := errorResponse.MarshalJson()
		return nil, fiber.NewError(fiber.StatusConflict, string(errorResponseJson))
	}

	// Check if thing is referenced in any other things
	dbCtx, dbSpan = otel.Tracer("prisma").Start(request.Context, "Check if thing is referenced in another thing")
	things, err := prisma.Client.Thing.FindMany(
		db.Thing.SpaceID.Equals(request.Input.SpaceId),
	).Select(
		db.Thing.Attributes.Field(),
	).Exec(dbCtx)
	if err != nil && err.Error() != "ErrNotFound" {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error checking if thing is used in another thing: %s", err)),
		)
		dbSpan.End()

		return nil, err
	}

	dbSpan.End()

	thingsWithThing := []db.ThingModel{}
	for _, thing := range things {
		var thingAttributes map[string]any
		err = json.Unmarshal(thing.Attributes, &thingAttributes)
		if err != nil {
			request.Span.SetAttributes(
				attribute.String("error.type", "parse_json"),
				attribute.String("error.message", fmt.Sprintf("could not read thing attributes: %s", err.Error())),
			)

			errorResponse := ioteahttp.NewErrorResponse([]any{
				fmt.Sprintf("could not read thing attributes: %s", err.Error()),
			})
			errorResponseJson, _ := errorResponse.MarshalJson()
			return nil, fiber.NewError(fiber.StatusBadRequest, string(errorResponseJson))
		}

		for _, value := range thingAttributes {
			if value == request.Input.ThingId {
				thingsWithThing = append(thingsWithThing, thing)
			}
		}
	}

	if len(thingsWithThing) > 0 {
		request.Span.SetAttributes(
			attribute.String("error.type", "context_validation"),
			attribute.String("error.message", fmt.Sprintf("thing is referenced in %d other thing(s)", len(thingsWithThing))),
		)

		errorResponse := ioteahttp.NewErrorResponse([]any{
			fmt.Sprintf("thing is referenced in %d other thing(s)", len(thingsWithThing)),
		})
		errorResponseJson, _ := errorResponse.MarshalJson()
		return nil, fiber.NewError(fiber.StatusConflict, string(errorResponseJson))
	}

	now := util.GetCurrentTime()

	attributes, err := json.Marshal(request.Input.Attributes)
	if err != nil {
		errMessage := fmt.Sprintf("error saving attributes: %s", err)
		request.Span.SetAttributes(
			attribute.String("error.type", "parse_json"),
			attribute.String("error.message", errMessage),
		)

		response := ioteahttp.NewErrorResponse([]any{errMessage})
		responseJson, _ := response.MarshalJson()
		return nil, fiber.NewError(400, string(responseJson))
	}

	dbCtx, dbSpan = otel.Tracer("prisma").Start(request.Context, "Update thing in space")
	thing, err := prisma.Client.Thing.FindUnique(
		db.Thing.ID.Equals(request.Input.ThingId),
	).Update(
		db.Thing.Name.Set(request.Input.Name),
		db.Thing.Attributes.Set(attributes),

		// TODO: add user ID to input and set updatedBy
		db.Thing.UpdatedBy.Set(request.GetActorId()),
		db.Thing.UpdatedAt.Set(now),
	).Exec(dbCtx)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error updating thing in the database: %s", err)),
		)
		dbSpan.End()

		return nil, err
	}
	dbSpan.End()

	output := Output{
		Thing: thing,
	}

	return &output, nil
}
