package api

import (
	"context"

	pb "github.com/ongruent/gruent/libs/protocols/controller"
)

type controllerGRPC struct {
	pb.UnimplementedControllerServiceServer
}

func newControllerGRPC() *controllerGRPC {
	return &controllerGRPC{}
}

func (*controllerGRPC) Publish(context.Context, *pb.PublishRequest) (*pb.PublishResponse, error) {
	return &pb.PublishResponse{}, nil
}

func (*controllerGRPC) Unpublish(context.Context, *pb.UnpublishRequest) (*pb.UnpublishResponse, error) {
	return &pb.UnpublishResponse{}, nil
}

func (*controllerGRPC) Status(context.Context, *pb.StatusRequest) (*pb.StatusResponse, error) {
	return &pb.StatusResponse{}, nil
}
