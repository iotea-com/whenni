package environmentsDelete

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	pb "github.com/iotea-com/iotea/libs/protocols/devenv"
	"github.com/iotea-com/iotea/prisma/db"
	devenvService "github.com/iotea-com/iotea/services/http-api/services/devenv"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// Assuming DevEnvServiceClient is your generated client for the service
type DevEnvServiceClient interface {
	DeleteEnvironment(ctx context.Context, in *pb.DeleteRequest, opts ...grpc.CallOption) (*db.DevEnvironmentModel, error)
}

func execute(request *ioteahttp.Request[Input]) (*Output, error) {

	tracer := otel.Tracer("devenv")
	_, span := tracer.Start(request.Context, "execute")
	defer span.End()
	span.AddEvent("Starting execution")

	// Prepare request
	grpcRequest := &pb.DeleteRequest{
		EnvironmentID: request.Input.EnvironmentID,
		Region:        request.Input.Region,
		SpaceId:       request.Input.SpaceId,
	}

	// Contacting the gRPC service
	if devenvService.Client == nil {
		errMsg := "DevEnv API is currently not reachable"
		span.AddEvent(errMsg)
		return nil, fiber.NewError(500, errMsg)
	}

	devenvServiceClient := *devenvService.DevEnvServiceClient
	response, err := devenvServiceClient.Delete(request.Context, grpcRequest)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			span.AddEvent(fmt.Sprintf("internal server error: %v", err))
			return nil, fiber.NewError(http.StatusInternalServerError, st.Message())
		}
		span.AddEvent(fmt.Sprintf("error contacting gRPC service: %v", err))
		return nil, err
	}

	span.AddEvent("successfully deleted environment")

	return &Output{
		Message: response.Message,
	}, nil
}
