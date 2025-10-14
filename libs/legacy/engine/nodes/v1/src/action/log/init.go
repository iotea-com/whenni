package main

import (
	"fmt"
	"time"

	"github.com/agoda-com/opentelemetry-go/otelzerolog"
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
)

func (n *LogActionNode) Init(params node.InitParams) node.Error {
	startTime := time.Now()

	// Setup tracer, meter and logger
	if err := n.setupObservability(params.Obsv); err != nil {
		return node.Error{
			Type:   node.FatalError,
			Reason: fmt.Sprintf("error setting up observability: %v", err),
		}
	}

	// Set notify channel
	n.notifyChannel = params.NotifyChannel

	n.logger.Debug().Msgf("%s Node Initialized OK in %vms", NodeName, time.Since(startTime).Milliseconds())

	return node.Error{
		Type: node.NoError,
	}
}

// setupObservability will create a tracer, a meter and a logger to be used
// internally by the node, based on the providers passed at initialization
func (n *LogActionNode) setupObservability(obsv node.Observability) error {
	n.tracer = obsv.TracerProvider.Tracer("LogActionNodeTracer")
	n.meter = obsv.MeterProvider.Meter("LogActionNodeMeter")

	hook := otelzerolog.NewHook(obsv.LoggerProvider)
	n.logger = obsv.Logger.With().Str("node", NodeName).Logger().Hook(hook)

	return nil
}
