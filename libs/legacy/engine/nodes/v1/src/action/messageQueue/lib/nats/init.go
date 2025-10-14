package natsActionNode

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

func (n *NatsActionSubnode) Init(params node.InitParams) node.Error {
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

	// Setup configuration fields
	if err := n.setConfig(params.Config); err != nil {
		return node.Error{
			Type:   node.ValidationError,
			Reason: fmt.Sprintf("invalid configuration, %v", err),
		}
	}

	// Connect to the Kafka cluster
	if err := n.connectToServer(); err != nil {
		return node.Error{
			Type:   node.FatalError,
			Reason: fmt.Sprintf("error connecting to NATS server: %v", err),
		}
	}

	n.logger.Debug().Msgf("%s Node Initialized OK in %vms", NodeName, time.Since(startTime).Milliseconds())

	return node.Error{
		Type: node.NoError,
	}
}

// setConfig unmarshals the passed config byte array into the node's configuration
// struct, and validates all the configuration fields
func (n *NatsActionSubnode) setConfig(config []byte) error {
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

// connectToServer will instantiate a new nats.Conn with options based on the user's
// configuration, and then attempt to connect to the provided cluster.
func (n *NatsActionSubnode) connectToServer() error {
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

	n.logger.Debug().Msgf("NATS client connected to %s:%d and topic %s is ready.", server.Host, server.Port, n.config.Topic)

	return nil
}

// setupObservability will create a tracer, a meter and a logger to be used
// internally by the node, based on the providers passed at initialization
func (n *NatsActionSubnode) setupObservability(obsv node.Observability) error {
	n.tracer = obsv.TracerProvider.Tracer("MessageQueueActionNodeTracer")
	n.meter = obsv.MeterProvider.Meter("MessageQueueActionNodeMeter")

	hook := otelzerolog.NewHook(obsv.LoggerProvider)
	n.logger = obsv.Logger.With().Str("node", NodeName).Logger().Hook(hook)

	return nil
}
