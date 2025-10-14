package natsSourceNode

import (
	"fmt"
	"time"

	"github.com/goccy/go-json"

	"github.com/agoda-com/opentelemetry-go/otelzerolog"
	"github.com/go-playground/validator/v10"
	"github.com/iotea-com/iotea/libs/engine/dependencies/things"
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	"github.com/nats-io/nats.go"
)

const (
	// The connection to a cluster should not take a very long time
	ClusterConnectTimeout = 1 * time.Second
)

func (n *NatsSourceSubnode) Init(params node.InitParams) node.Error {
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

	// Setup the error channel
	n.errChan = make(chan error)

	// Setup configuration fields
	if err := n.setConfig(params.Config); err != nil {
		return node.Error{
			Type:   node.ValidationError,
			Reason: fmt.Sprintf("invalid configuration, %v", err),
		}
	}

	// Connect to the NATS server
	err := n.connectToServer()
	if err != nil {
		return node.Error{
			Type:   node.FatalError,
			Reason: err.Error(),
		}
	}

	n.logger.Debug().Msgf("%s Node Initialized OK in %vms", NodeName, time.Since(startTime).Milliseconds())

	return node.Error{
		Type: node.NoError,
	}
}

// setConfig unmarshals the passed config byte array into the node's configuration
// struct, and validates all the configuration fields
func (n *NatsSourceSubnode) setConfig(config []byte) error {
	if err := json.Unmarshal(config, &n.config); err != nil {
		return fmt.Errorf("invalid node configuration JSON passed: %v", err)
	}

	// Validate config
	v := validator.New()

	if err := v.Struct(&n.config); err != nil {
		return err
	}

	return nil
}

// connectToServer will instantiate a new nats.Conn with options based on the user's
// configuration, and then attempt to connect to the provided server.
func (n *NatsSourceSubnode) connectToServer() error {
	// Get a things.NatsServer from config
	jsonServer, err := json.Marshal(n.config.ThingNatsClient.Server)
	if err != nil {
		return fmt.Errorf("error marshalling cluster: %s", err)
	}

	var server things.NatsServer
	err = json.Unmarshal(jsonServer, &server)
	if err != nil {
		return fmt.Errorf("error unmarshalling JSON into struct: %s", err)
	}

	// Connect to the server
	conn, err := nats.Connect(fmt.Sprintf("nats://%s:%d", server.Host, server.Port))
	if err != nil {
		return fmt.Errorf("error connecting to NATS server: %v", err)
	}

	n.conn = conn

	// Subscribe to messages
	if !n.conn.IsConnected() {
		return fmt.Errorf("cannot subscribe to the provided topic on the NATS server due to an inactive NATS connection")
	}

	_, err = n.conn.Subscribe(n.config.Topic, n.natsSubscriptionCallback)
	if err != nil {
		return fmt.Errorf("error subscribing to NATS server: %s", err)
	}

	n.logger.Debug().Msgf("NATS client connected and client subscribed to topic %s.", n.config.Topic)

	return nil
}

// setupObservability will create a tracer, a meter and a logger to be used
// internally by the node, based on the providers passed at initialization
func (n *NatsSourceSubnode) setupObservability(obsv node.Observability) error {
	n.tracer = obsv.TracerProvider.Tracer("MessageQueueSourceNodeTracer")
	n.meter = obsv.MeterProvider.Meter("MessageQueueSourceNodeMeter")

	hook := otelzerolog.NewHook(obsv.LoggerProvider)
	n.logger = obsv.Logger.With().Str("node", NodeName).Logger().Hook(hook)

	return nil
}
