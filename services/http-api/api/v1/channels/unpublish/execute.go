package channelsUnpublish

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

	// Send Unpublish message to the orchestrator
	ctx, cancel := context.WithTimeout(request.Context, time.Second)
	defer cancel()
	message := &channelService.TerminateRequest{ChannelId: request.Input.ChannelId}
	terminateReponse, err := client.Terminate(ctx, message)
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "grpc_request"),
			attribute.String("error.message", fmt.Sprintf("could not send unpublish message to the orchestrator: %s", err)),
		)
		response := ioteahttp.NewErrorResponse([]any{
			"Could not unpublish the channel at this time. Please try again.",
		})
		responseJson, _ := response.MarshalJson()
		return nil, fiber.NewError(fiber.StatusServiceUnavailable, string(responseJson))
	} else if terminateReponse.Error != "" {
		request.Span.SetAttributes(
			attribute.String("error.type", "grpc_request"),
			attribute.String("error.message", fmt.Sprintf("an unpublish error ocurred: %s", terminateReponse.Error)),
		)
		response := ioteahttp.NewErrorResponse([]any{
			fmt.Sprintf("Could not unpublish the channel at this time: %s, Please try again.", terminateReponse.Error),
		})
		responseJson, _ := response.MarshalJson()
		return nil, fiber.NewError(fiber.StatusBadRequest, string(responseJson))
	}

	return nil, nil
}
