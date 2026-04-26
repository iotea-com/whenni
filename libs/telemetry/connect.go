package telemetry

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

func connectToCollector(ctx context.Context, config Config) (*grpc.ClientConn, error) {
	// Determine the gRPC dial options based on the URL scheme
	var dialOpts []grpc.DialOption
	// if config.UseTLS {
	// 	// Use Journey CA certificate for TLS connections
	// 	caCert, _ := base64.StdEncoding.DecodeString(env.Get("JOURNEY_CA_CERT"))
	// 	certPool := x509.NewCertPool()
	// 	certPool.AppendCertsFromPEM(caCert)
	// 	creds := credentials.NewClientTLSFromCert(certPool, "")
	// 	dialOpts = append(dialOpts, grpc.WithTransportCredentials(creds))
	// } else {
	// Use insecure connection for HTTP
	dialOpts = append(dialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	// }

	// fmt.Printf("Connecting to collector at %s with UseTLS = %t\n", config.OtelCollectorEndpoint, config.UseTLS)
	fmt.Printf("Connecting to collector at %s with UseTLS = %t\n", config.OtelCollectorEndpoint, false)

	// Set up gRPC connection to the collector
	conn, err := grpc.NewClient(config.OtelCollectorEndpoint, dialOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %s", err)
	}

	// Wait for the connection to be ready
	conn.Connect() // Trigger connection attempt
	for {
		time.Sleep(100 * time.Millisecond)
		if conn == nil {
			return nil, fmt.Errorf("failed to create gRPC connection to collector: connection is nil")
		}
		state := conn.GetState()
		if state == connectivity.Ready {
			fmt.Println("Successfully connected to collector")
			break
		}
		if state == connectivity.TransientFailure || state == connectivity.Shutdown {
			return nil, fmt.Errorf("failed to connect to collector at %v: connection state is %v", config.OtelCollectorEndpoint, state)
		}
		if !conn.WaitForStateChange(ctx, state) {
			return nil, fmt.Errorf("timeout waiting for connection to collector at %v", config.OtelCollectorEndpoint)
		}
	}

	return conn, nil
}
