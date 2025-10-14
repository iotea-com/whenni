package channelsPublish

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/iotea-com/iotea/libs/engine/channels"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	channelService "github.com/iotea-com/iotea/libs/protocols/channels"
	pbChannel "github.com/iotea-com/iotea/libs/protocols/channels"
	"github.com/iotea-com/iotea/prisma/db"
	"github.com/iotea-com/iotea/services/http-api/config"
	"github.com/iotea-com/iotea/services/http-api/services/prisma"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func execute(request *ioteahttp.Request[Input]) (*Output, error) {
	request.Span.AddEvent("execute")
	request.Span.SetAttributes(
		attribute.String("request.Input.ChannelId", request.Input.ChannelId),
		attribute.String("request.Input.SpaceId", request.Input.SpaceId),
	)

	// Set up rules engine client
	// TODO: Add TLS to this connection for production
	engineConn, err := grpc.NewClient(config.VaultConf.EngineGrpcServiceUrl,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)

	defer func() {
		if err := engineConn.Close(); err != nil {
			request.Span.SetAttributes(
				attribute.String("error.type", "grpc_request"),
				attribute.String("error.message", fmt.Sprintf("failed to close rules engine grpc client: %v", err)),
			)
		}
	}()

	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "grpc_request"),
			attribute.String("error.message", fmt.Sprintf("failed to setup rules engine grpc client: endpoint is %v and error is %v", config.VaultConf.EngineGrpcServiceUrl, err)),
		)
		response := ioteahttp.NewErrorResponse([]any{
			"Could not get the status of the channel at this time. Please try again.",
		})
		responseJson, _ := response.MarshalJson()
		return nil, fiber.NewError(fiber.StatusServiceUnavailable, string(responseJson))
	}

	client := pbChannel.NewChannelServiceClient(engineConn)

	// Send Publish message to the orchestrator
	ctx, cancel := context.WithTimeout(request.Context, time.Second*5)
	defer cancel()
	message := &channelService.ExecuteRequest{ChannelId: request.Input.ChannelId}
	executeResponse, err := client.Execute(ctx, message)
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "grpc_request"),
			attribute.String("error.message", fmt.Sprintf("could not send execute response to the orchestrator: %s", err)),
		)
		response := ioteahttp.NewErrorResponse([]any{
			"Could not publish the channel at this time. Please try again.",
		})
		responseJson, _ := response.MarshalJson()
		return nil, fiber.NewError(fiber.StatusServiceUnavailable, string(responseJson))
	} else if executeResponse.Error != "" {
		request.Span.SetAttributes(
			attribute.String("error.type", "grpc_request"),
			attribute.String("error.message", fmt.Sprintf("a publish error ocurred: %s", executeResponse.Error)),
		)
		response := ioteahttp.NewErrorResponse([]any{
			fmt.Sprintf("Could not publish the channel at this time: %s, Please try again.", executeResponse.Error),
		})
		responseJson, _ := response.MarshalJson()
		return nil, fiber.NewError(fiber.StatusBadRequest, string(responseJson))
	}

	return nil, nil
}

func isSmallRuntime(channelId string) (bool, error) {
	// Fetch channel from database
	channelJson, err := prisma.Client.Channel.FindUnique(
		db.Channel.ID.Equals(channelId),
	).Select(
		db.Channel.ID.Field(),
		db.Channel.Config.Field(),
		db.Channel.SpaceID.Field(),
	).Exec(context.Background())

	if err != nil {
		return false, fmt.Errorf("error fetching channel from database: %s", err)
	}

	// Unmarshal channel config
	var c channels.Channel
	if err := json.Unmarshal(channelJson.Config, &c); err != nil {
		return false, fmt.Errorf("error unmarshalling channel: %s", err)
	}

	// Check if the channel is a small runtime
	if c.Runtime.Size != "small" {
		return false, nil
	}

	return true, nil
}
