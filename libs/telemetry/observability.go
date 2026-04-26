package telemetry

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/iotea-com/iotea/libs/legacy/engine/environment"
	"github.com/iotea-com/iotea/services/http-api/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

// Initializes all otel providers and returns a struct with references to the providers and a
// shutdown function to properly flush and close resources.
func InitObservability(serviceName string) (*Observability, error) {
	serviceInstance := (func() string {
		if config.Env == environment.Development {
			return "docker-compose"
		}

		return uuid.NewString()
	})()
	obsvConfig := Config{
		OtelCollectorEndpoint: config.VaultConf.OtelCollectorEndpoint,
		ResourceAttr: ResourceAttributes{
			Environment:     config.Env.String(),
			ServiceName:     serviceName,
			ServiceInstance: serviceInstance,
		},
		UseTLS: config.Env != environment.Development,
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
