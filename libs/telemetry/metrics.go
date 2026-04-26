package telemetry

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"google.golang.org/grpc"
)

var (
	meter                     = otel.Meter("api")
	ExecutionSubmitCounter, _ = meter.Int64Counter("execution.submit.count")
	SessionSubmitCounter, _   = meter.Int64Counter("session.submit.count")
	ExecutionIDCounter, _     = meter.Int64Counter("execution.id.count")
	UserIDCounter, _          = meter.Int64Counter("user.id.count")
	PipelineIDCounter, _      = meter.Int64Counter("pipeline.id.count")
	PipelineStageCounter, _   = meter.Int64Counter("pipeline.stage.count")
)

// Sets up the meter provider. If no otel collector connection is provided, a noop meter provider
// will be used.
func setupMeterProvider(ctx context.Context, conn *grpc.ClientConn, res *resource.Resource) (*metric.MeterProvider, error) {
	if conn == nil {
		mp := metric.NewMeterProvider()
		return mp, nil
	}

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
