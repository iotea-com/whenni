package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/iotea-com/iotea/libs/telemetry"
	"github.com/iotea-com/iotea/services/http-api/api"
	"github.com/iotea-com/iotea/services/http-api/config"
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
	// Setup observability
	observability, err := telemetry.InitObservability(telemetry.InitConfig{
		ServiceName:           "http-api",
		OtelCollectorEndpoint: config.VaultConf.OtelCollectorEndpoint,
		Environment:           config.Env,
	})
	if err != nil {
		panic(fmt.Errorf("failed to setup observability: %v", err))
	}
	defer observability.Shutdown()

	// Create context for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Create new API server
	api := api.New()

	// Start server in a goroutine
	go func() {
		if err := api.Listen(config.VaultConf.ApiPort); err != nil {
			panic(fmt.Errorf("Server error: %s", err))
		}
	}()

	// Wait for interrupt signal
	<-ctx.Done()
	// Only print the shutdown message once
	if os.Getppid() <= 1 {
		fmt.Println("Shutting down server...")
	}

	// Shut down the server
	errors := api.Shutdown()
	if len(errors) > 0 {
		fmt.Printf("Server shutdown errors for PID %v: %v\n", os.Getppid(), errors)
	}

	// Only print the shutdown message once
	fmt.Printf("Server gracefully stopped for PID %v\n", os.Getppid())
}
