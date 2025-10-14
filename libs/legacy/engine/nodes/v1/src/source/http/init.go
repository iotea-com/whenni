package main

import (
	"fmt"
	"time"

	"github.com/goccy/go-json"

	"github.com/agoda-com/opentelemetry-go/otelzerolog"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	router "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/source/http/lib"
	"github.com/iotea-com/iotea/libs/val"
)

func (n *HttpSourceNode) Init(params node.InitParams) node.Error {
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

	// Find the right input channel where the responses are sent
	// by the response nodes into this node
	var responsesInputChannel node.IoChannel
	for _, channel := range n.inputChannels {
		if channel.Id == ResponsesInputChannelId {
			responsesInputChannel = channel
		}
	}

	// Initialize the response router
	n.router = router.NewResponseRouter(n.logger)

	// Start listening on the responses input channel
	go func(channel node.IoChannel) {
		for data := range channel.Channel {
			if err := n.router.RouteResponse(channel.Id, data); err != nil {
				n.logger.Error().Msgf("Failed to route response: %v", err)
			}
		}
	}(responsesInputChannel)

	// Initialize
	app := fiber.New(fiber.Config{
		ReadTimeout:      10 * time.Second,
		WriteTimeout:     10 * time.Second,
		DisableKeepalive: true, // Important for graceful shutdown
	})

	handler := func(c *fiber.Ctx) error {
		return n.httpRequestSubscriptionCallback(c)
	}

	// Register the handler for all methods
	app.All("/", handler)

	n.server = app

	n.logger.Debug().Msgf("%s Node Initialized OK in %vms", NodeName, time.Since(startTime).Milliseconds())

	return node.Error{
		Type: node.NoError,
	}
}

// setConfig unmarshalls the passed config byte array into the node's configuration
// struct, and validates all the configuration fields
func (n *HttpSourceNode) setConfig(config []byte) error {
	if err := json.Unmarshal(config, &n.config); err != nil {
		return fmt.Errorf("invalid node configuration JSON passed: %v", err)
	}

	v := validator.New()
	v.RegisterValidation("valid_paths", val.IsValidPaths)

	if err := v.Struct(&n.config); err != nil {
		return err
	}

	// Create channel to communicate all errors created by go routines in this node
	n.errChan = make(chan error)

	return nil
}

// setupObservability will create a tracer, a meter and a logger to be used
// internally by the node, based on the providers passed at initialization
func (n *HttpSourceNode) setupObservability(obsv node.Observability) error {
	n.tracer = obsv.TracerProvider.Tracer("HttpSourceNodeTracer")
	n.meter = obsv.MeterProvider.Meter("HttpSourceNodeMeter")

	hook := otelzerolog.NewHook(obsv.LoggerProvider)
	n.logger = obsv.Logger.With().Str("node", NodeName).Logger().Hook(hook)

	return nil
}
