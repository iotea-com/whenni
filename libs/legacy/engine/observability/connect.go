package observability

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func connectToCollector(ctx context.Context, config Config) (*grpc.ClientConn, error) {
	// Determine the gRPC dial options based on the URL scheme
	var dialOpts []grpc.DialOption
	if strings.HasPrefix(config.OtelCollectorEndpoint, "https://") {
		// Use default system CA certificates for HTTPS connections
		creds := credentials.NewClientTLSFromCert(nil, "")
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(creds))
	} else {
		// Remove "http://" from the beginning of the connection string
		config.OtelCollectorEndpoint = strings.TrimPrefix(config.OtelCollectorEndpoint, "http://")

		// Use insecure connection for HTTP
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// Set up gRPC connection to the collector
	conn, err := grpc.NewClient(config.OtelCollectorEndpoint, dialOpts...)

	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector at %v: %w", config.OtelCollectorEndpoint, err)
	}

	return conn, nil
}
