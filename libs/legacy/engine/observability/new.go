package observability

import (
	"context"
	"fmt"

	"github.com/agoda-com/opentelemetry-logs-go/logs"
	"github.com/go-playground/validator/v10"
	"github.com/iotea-com/iotea/libs/engine/environment"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
)

type ResourceAttributes struct {
	// Environment (required)
	Environment environment.Env `validate:"required"`

	// Service (required)
	ServiceName     string `validate:"required"`
	ServiceInstance string `validate:"required"`

	// Build Information (required)
	BuildVersion string `validate:"required"`
	BuildTime    string `validate:"required"`
	BuildCommit  string `validate:"required"`
	BuildDirty   string `validate:"required"`
	BuildCreator string `validate:"required"`
}

type Config struct {
	// Otel end-point connection information
	OtelCollectorEndpoint string

	// Resource information
	ResourceAttr ResourceAttributes
}

// Obsv holds references to the providers initialized by the SetupObservability function
// and a shutdown function to properly close resources.
type Obsv struct {
	TracerProvider trace.TracerProvider
	MeterProvider  metric.MeterProvider
	LoggerProvider logs.LoggerProvider
	Shutdown       func(context.Context) error
}

func New(ctx context.Context, config Config) (*Obsv, error) {
	// Validate the config
	if err := validator.New().Struct(config); err != nil {
		return nil, fmt.Errorf("invalid resource attributes: %w", err)
	}

	// Set the attributes
	attributes := setAttributes(config.ResourceAttr.Environment, config.ResourceAttr)

	// Create the resource
	res, err := resource.New(ctx, resource.WithAttributes(attributes...))
	if err != nil {
		return nil, fmt.Errorf("could not create resource: %v", err)
	}

	conn, err := connectToCollector(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("could not connect to collector: %v", err)
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
	shutdown := func(ctx context.Context) error {
		// Force flush all providers
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

	return &Obsv{
		TracerProvider: tracerProvider,
		MeterProvider:  meterProvider,
		LoggerProvider: loggerProvider,
		Shutdown:       shutdown,
	}, nil
}
