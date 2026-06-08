package main

import (
	"fmt"
	"os"

	"github.com/iotea-com/iotea/libs/legacy/engine/environment"
	"github.com/iotea-com/iotea/libs/telemetry"
	"github.com/iotea-com/iotea/services/controller/api"
	"github.com/iotea-com/iotea/services/controller/config"
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
