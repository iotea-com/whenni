package devenvService

import (
	pbDevenv "github.com/iotea-com/iotea/libs/protocols/devenv"
	"google.golang.org/grpc"
)

// Rules engine grpc client
var Client *grpc.ClientConn

// Dev environment service client
var DevEnvServiceClient *pbDevenv.DevEnvClient
