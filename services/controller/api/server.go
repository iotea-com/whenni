package api

import (
	"fmt"
	"log"
	"net"

	pb "github.com/iotea-com/iotea/libs/protocols/controller"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

// Server wraps a gRPC server instance.
type Server struct {
	grpc *grpc.Server
}

// New constructs a basic gRPC server with default options.
func New() *Server {
	s := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)
	pb.RegisterControllerServiceServer(s, newControllerGRPC())
	return &Server{grpc: s}
}

// Listen binds to port and serves until the process exits or Serve returns an error.
func (s *Server) Listen(port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	if err := s.grpc.Serve(lis); err != nil {
		log.Fatalf("serve: %v", err)
	}
}
