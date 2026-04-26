package telemetry

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

type ResourceAttributes struct {
	// Environment (required)
	Environment string `validate:"required"`

	// Service (required)
	ServiceName     string `validate:"required"`
	ServiceInstance string `validate:"required"`
}

type Config struct {
	// Otel end-point connection information
	OtelCollectorEndpoint string
	UseTLS                bool

	// Resource information
	ResourceAttr ResourceAttributes
}

// Observability holds references to the providers initialized by the SetupObservability function
// and a shutdown function to properly close resources.
type Observability struct {
	TracerProvider trace.TracerProvider
	MeterProvider  metric.MeterProvider
	LoggerProvider *sdklog.LoggerProvider
	Shutdown       func() error
}

// Sets up the observability providers and returns a struct with references to the providers and a
// shutdown function to properly flush and close resources.
func New(ctx context.Context, config Config) (*Observability, error) {
	// Validate the config
	if err := validator.New().Struct(config); err != nil {
		return nil, fmt.Errorf("invalid resource attributes: %w", err)
	}

	res, err := resource.New(ctx, resource.WithAttributes(
		attribute.String("environment", config.ResourceAttr.Environment),
		attribute.String("service.name", config.ResourceAttr.ServiceName),
		attribute.String("service.instance", config.ResourceAttr.ServiceInstance),
	))
	if err != nil {
		return nil, fmt.Errorf("could not create resource: %v", err)
	}

	var conn *grpc.ClientConn = nil
	if config.OtelCollectorEndpoint != "" {
		conn, err = connectToCollector(ctx, config)
		if err != nil {
			fmt.Printf("Could not connect to collector: %v\n - only stdout logging will be available; traces and metrics are disabled\n", err)
		}
	} else {
		fmt.Println("No otel collector endpoint was provided - only stdout logging will be available; traces and metrics are disabled")
	}

	// Initialize TracerProvider
	tracerProvider, err := setupTracerProvider(ctx, conn, res)
	if err != nil {
		return nil, err
	}

	// Initialize MeterProvider
	meterProvider, err := setupMeterProvider(ctx, conn, res)
	if err != nil {
		return nil, err
	}

	// Initialize LoggerProvider
	loggerProvider, err := setupLoggerProvider(ctx, conn, res)
	if err != nil {
		return nil, err
	}

	// Define a shutdown function to gracefully close resources
	shutdown := func() error {
		// Force flush all providers
		ctx := context.Background()
		if err := tracerProvider.ForceFlush(ctx); err != nil {
			return err
		}
		if err := meterProvider.ForceFlush(ctx); err != nil {
			return err
		}
		if err := loggerProvider.ForceFlush(ctx); err != nil {
			return err
		}

		// Shutdown all providers
		if err := tracerProvider.Shutdown(ctx); err != nil {
			return err
		}
		if err := meterProvider.Shutdown(ctx); err != nil {
			return err
		}
		if err := loggerProvider.Shutdown(ctx); err != nil {
			return err
		}

		return nil
	}

	return &Observability{
		TracerProvider: tracerProvider,
		MeterProvider:  meterProvider,
		LoggerProvider: loggerProvider,
		Shutdown:       shutdown,
	}, nil
}
