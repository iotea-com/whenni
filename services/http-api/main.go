package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/iotea-com/iotea/libs/legacy/engine/observability"
	pbDevenv "github.com/iotea-com/iotea/libs/protocols/devenv"
	"github.com/iotea-com/iotea/services/http-api/api"
	"github.com/iotea-com/iotea/services/http-api/config"
	devenvService "github.com/iotea-com/iotea/services/http-api/services/devenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

// @title IOTEA
// @version alpha-v1.0
// @description Bringing connectivity and automation to everyone, everywhere.
// @schemes https http
// @host localhost:9001
// @BasePath /v1
// @securitydefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	// Setup OTLP trace exporter
	providers := initObservability()
	defer providers.Shutdown(context.Background())

	// Create context for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Set up dev environment service client
	// TODO: Add TLS to this connection for production
	devenvConn, err := grpc.NewClient(config.VaultConf.DevenvGrpcServerUrl,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		panic(fmt.Errorf("failed to setup devenv grpc client: endpoint is %v and error is %v", config.VaultConf.DevenvGrpcServerUrl, err))
	}
	// Ensure connection is closed on exit
	defer devenvConn.Close()

	devenvServiceClient := pbDevenv.NewDevEnvClient(devenvConn)
	devenvService.Client = devenvConn
	devenvService.DevEnvServiceClient = &devenvServiceClient

	// Create new API server
	api := api.New()

	// Start server in a goroutine
	go func() {
		if err := api.Listen(config.VaultConf.ApiPort); err != nil {
			log.Fatal().Msgf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-ctx.Done()
	// Only print the shutdown message once
	if os.Getppid() <= 1 {
		log.Info().Msg("Shutting down server...")
	}

	// Shut down the server
	errors := api.Shutdown()
	if len(errors) > 0 {
		log.Error().Msgf("Server shutdown errors for PID %v: %v", os.Getppid(), errors)
	}

	// Only print the shutdown message once
	log.Info().Msgf("Server gracefully stopped for PID %v", os.Getppid())
}

func initObservability() *observability.Obsv {
	obsvConfig := observability.Config{
		OtelCollectorEndpoint: config.VaultConf.PlatformOtelCollectorEndpoint,
		ResourceAttr: observability.ResourceAttributes{
			Environment:     config.Env,
			ServiceName:     "HTTP API",
			ServiceInstance: "1.0.0", //TODO: get instance id
			BuildVersion:    config.BuildVersion,
			BuildTime:       config.BuildTime,
			BuildCommit:     config.BuildCommit,
			BuildDirty:      config.BuildDirty,
			BuildCreator:    config.BuildCreator,
		},
	}

	setupCtx, cancelFunc := context.WithTimeout(context.Background(), time.Millisecond*100)
	providers, err := observability.New(setupCtx, obsvConfig)
	if err != nil {
		panic(fmt.Errorf("failed to setup observability: endpoint is %v and error is %v", config.VaultConf.PlatformOtelCollectorEndpoint, err))
	}
	cancelFunc()

	// Set global providers
	otel.SetTracerProvider(providers.TracerProvider)
	otel.SetMeterProvider(providers.MeterProvider)

	// Set global propagator to tracecontext
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// Set global logging
	// hook := otelzerolog.NewHook(providers.LoggerProvider)
	// log.Logger = log.Hook(hook)

	return providers
}
