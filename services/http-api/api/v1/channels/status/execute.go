package channelsStatus

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	channelService "github.com/iotea-com/iotea/libs/protocols/channels"
	pbChannel "github.com/iotea-com/iotea/libs/protocols/channels"
	"github.com/iotea-com/iotea/services/http-api/config"
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

	// Send message to channel service
	ctx, cancel := context.WithTimeout(request.Context, time.Second)
	defer cancel()
	message := &channelService.StatusRequest{ChannelId: request.Input.ChannelId}
	response, err := client.Status(ctx, message)
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "grpc_request"),
			attribute.String("error.message", fmt.Sprintf("could not connect to the orchestrator: %s", err)),
		)
		response := ioteahttp.NewErrorResponse([]any{
			"Could not get the status of the channel at this time. Please try again.",
		})
		responseJson, _ := response.MarshalJson()
		return nil, fiber.NewError(fiber.StatusServiceUnavailable, string(responseJson))
	}

	output := &Output{
		Status: response.Status,
	}

	return output, nil
}
