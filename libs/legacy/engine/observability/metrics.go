package observability

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"google.golang.org/grpc"
)

func setupMeterProvider(ctx context.Context, conn *grpc.ClientConn, res *resource.Resource) (*metric.MeterProvider, error) {
	// Initialize the OTLP metric exporter using gRPC
	metricExporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create metric exporter: %w", err)
	}

	// Configure the MeterProvider with a 1-second interval for the periodic reader
	pusher := metric.NewPeriodicReader(metricExporter,
		// Set the export interval to 1 second
		metric.WithInterval(1*time.Second),
	)

	// Create the MeterProvider with the exporter and resource
	mp := metric.NewMeterProvider(
		metric.WithReader(pusher),
		metric.WithResource(res),
	)

	return mp, nil
}
