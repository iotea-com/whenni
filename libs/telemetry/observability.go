package telemetry

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ongruent/gruent/libs/legacy/engine/environment"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type InitConfig struct {
	ServiceName           string
	OtelCollectorEndpoint string
	Environment           environment.Env
}

// Initializes all otel providers and returns a struct with references to the providers and a
// shutdown function to properly flush and close resources.
func InitObservability(cfg InitConfig) (*Observability, error) {
	serviceInstance := (func() string {
		if cfg.Environment == environment.Development {
			return "local"
		}

		return uuid.NewString()
	})()
	obsvConfig := Config{
		OtelCollectorEndpoint: cfg.OtelCollectorEndpoint,
		ResourceAttr: ResourceAttributes{
			Environment:     cfg.Environment.String(),
			ServiceName:     cfg.ServiceName,
			ServiceInstance: serviceInstance,
		},
		UseTLS: cfg.Environment != environment.Development,
	}

	setupCtx, cancelFunc := context.WithTimeout(context.Background(), time.Second)
	providers, err := New(setupCtx, obsvConfig)
	if err != nil {
		cancelFunc()
		return nil, err
	}
	cancelFunc()

	// Set global providers
	otel.SetTracerProvider(providers.TracerProvider)
	otel.SetMeterProvider(providers.MeterProvider)

	// Set global propagator to tracecontext
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return providers, nil
}
