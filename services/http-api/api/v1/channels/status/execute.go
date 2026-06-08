package channelsStatus

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	pbController "github.com/iotea-com/iotea/libs/protocols/controller"
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

	// Set up controller client
	// TODO: Add TLS to this connection for production
	controllerConn, err := grpc.NewClient(config.VaultConf.ControllerGrpcServiceUrl,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)

	defer func() {
		if err := controllerConn.Close(); err != nil {
			request.Span.SetAttributes(
				attribute.String("error.type", "grpc_request"),
				attribute.String("error.message", fmt.Sprintf("failed to close controller grpc client: %v", err)),
			)
		}
	}()

	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "grpc_request"),
			attribute.String("error.message", fmt.Sprintf("failed to setup controller grpc client: endpoint is %v and error is %v", config.VaultConf.ControllerGrpcServiceUrl, err)),
		)
		response := ioteahttp.NewErrorResponse([]any{
			"Could not get the status of the channel at this time. Please try again.",
		})
		responseJson, _ := response.MarshalJson()
		return nil, fiber.NewError(fiber.StatusServiceUnavailable, string(responseJson))
	}

	client := pbController.NewControllerServiceClient(controllerConn)

	// Send message to controller
	ctx, cancel := context.WithTimeout(request.Context, time.Second)
	defer cancel()
	statusResponse, err := client.Status(ctx, &pbController.StatusRequest{
		ChannelId: request.Input.ChannelId,
	})
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "grpc_request"),
			attribute.String("error.message", fmt.Sprintf("could not get the status of the channel from the controller: %s", err)),
		)
		response := ioteahttp.NewErrorResponse([]any{
			"Could not get the status of the channel at this time. Please try again.",
		})
		responseJson, _ := response.MarshalJson()
		return nil, fiber.NewError(fiber.StatusServiceUnavailable, string(responseJson))
	}

	output := &Output{
		Status: statusResponse.GetStatus().String(),
	}

	return output, nil
}
