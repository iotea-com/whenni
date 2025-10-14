package main

import (
	"fmt"
	"time"

	"github.com/goccy/go-json"

	"github.com/agoda-com/opentelemetry-go/otelzerolog"
	"github.com/go-playground/validator/v10"
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	"github.com/robfig/cron/v3"
)

const (
	// This timeout should allow enough room for a TLS handshake
	ConnectionTimeout = 5 * time.Second
)

func (n *TimerSourceNode) Init(params node.InitParams) node.Error {
	startTime := time.Now()

	// Setup tracer, meter and logger
	if err := n.setupObservability(params.Obsv); err != nil {
		return node.Error{
			Type:   node.FatalError,
			Reason: err.Error(),
		}
	}

	// Set the notify channel
	n.notifyChannel = params.NotifyChannel

	// Setup configuration fields
	if err := n.setConfig(params.Config); err != nil {
		return node.Error{
			Type:   node.ValidationError,
			Reason: fmt.Sprintf("invalid configuration, %v", err),
		}
	}

	// Add a new channel to the node config for sending trigger messages
	n.triggerChannel = make(chan bool, 1)

	// Add a new CRON scheduler to the node config
	n.cronScheduler = cron.New(cron.WithSeconds())

	// Schedule the task using the provided CRON expression
	_, err := n.cronScheduler.AddFunc(n.config.CronExpression, func() {
		n.triggerChannel <- true
	})
	if err != nil {
		return node.Error{
			Type:   node.FatalError,
			Reason: fmt.Sprintf("error starting timer schedule: %s", err),
		}
	}

	n.logger.Debug().Msgf("%s Node Initialized OK in %vms", NodeName, time.Since(startTime).Milliseconds())

	return node.Error{
		Type: node.NoError,
	}
}

// setConfig unmarshalls the passed config byte array into the node's configuration
// struct, and validates all the configuration fields
func (n *TimerSourceNode) setConfig(config []byte) error {
	if err := json.Unmarshal(config, &n.config); err != nil {
		return fmt.Errorf("invalid node configuration JSON passed: %v", err)
	}

	// Validate config
	v := validator.New()

	if err := v.Struct(n.config); err != nil {
		return err
	}

	return nil
}

// setupObservability will create a tracer, a meter and a logger to be used
// internally by the node, based on the providers passed at initialization
func (n *TimerSourceNode) setupObservability(obsv node.Observability) error {
	n.tracer = obsv.TracerProvider.Tracer("TriggerSourceNodeTracer")
	n.meter = obsv.MeterProvider.Meter("TriggerSourceNodeMeter")

	hook := otelzerolog.NewHook(obsv.LoggerProvider)
	n.logger = obsv.Logger.With().Str("node", NodeName).Logger().Hook(hook)

	return nil
}
