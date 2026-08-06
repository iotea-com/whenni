package main

import (
	"fmt"
	"os"

	"github.com/ongruent/gruent/libs/legacy/engine/environment"
	"github.com/ongruent/gruent/libs/telemetry"
	"github.com/ongruent/gruent/services/controller/api"
	"github.com/ongruent/gruent/services/controller/config"
)

func main() {
	observability, err := telemetry.InitObservability(telemetry.InitConfig{
		ServiceName:           "controller",
		OtelCollectorEndpoint: os.Getenv("OTEL_COLLECTOR_ENDPOINT"),
		Environment:           environment.GetFromEnvVar(),
	})
	if err != nil {
		panic(fmt.Errorf("failed to setup observability: %v", err))
	}
	defer observability.Shutdown()

	server := api.New()
	fmt.Println("Starting controller server on port 9002...")
	server.Listen(config.VaultConf.ControllerServerPort)
}
