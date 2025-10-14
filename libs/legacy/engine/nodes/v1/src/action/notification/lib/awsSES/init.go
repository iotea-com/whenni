package awsSES

import (
	"fmt"
	"time"

	"github.com/goccy/go-json"

	"github.com/agoda-com/opentelemetry-go/otelzerolog"
	"github.com/go-playground/validator/v10"
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	"github.com/iotea-com/iotea/libs/engine/template"
)

func (n *AwsSESActionSubnode) Init(params node.InitParams) node.Error {
	startTime := time.Now()

	// Setup tracer, meter and logger
	if err := n.setupObservability(params.Obsv); err != nil {
		return node.Error{
			Type:   node.ValidationError,
			Reason: err.Error(),
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

	n.logger.Debug().Msgf("%s Node Initialized OK in %vms", NodeName, time.Since(startTime).Milliseconds())

	return node.Error{
		Type: node.NoError,
	}
}

// setConfig unmarshals the passed config byte array into the node's configuration
// struct, and validates all the configuration fields
func (n *AwsSESActionSubnode) setConfig(config []byte) error {
	if err := json.Unmarshal(config, &n.config); err != nil {
		return fmt.Errorf("invalid node configuration JSON passed: %v", err)
	}

	// Validate config
	v := validator.New()
	// v.RegisterValidation("valid_paths", val.IsValidPaths)

	if err := v.Struct(&n.config); err != nil {
		return err
	}

	// Validate and setup the input schema - input schema is optional
	if n.config.ModelInput != nil {
		if err := n.config.ModelInput.Validate(); err != nil {
			return fmt.Errorf("invalid model input schema: %v", err)
		}

		// Create & verify template
		n.template = template.Template{
			Value:         n.config.TemplateMessage,
			ModelInput:    n.config.ModelInput,
			DefaultValues: n.config.DefaultValues,
		}

		if err := n.template.Verify(); err != nil {
			return err
		}
	}

	return nil
}

// setupObservability will create a tracer, a meter and a logger to be used
// internally by the node, based on the providers passed at initialization
func (n *AwsSESActionSubnode) setupObservability(obsv node.Observability) error {
	n.tracer = obsv.TracerProvider.Tracer("NotificationActionNodeTracer")
	n.meter = obsv.MeterProvider.Meter("NotificationActionNodeMeter")

	hook := otelzerolog.NewHook(obsv.LoggerProvider)
	n.logger = obsv.Logger.With().Str("node", NodeName).Logger().Hook(hook)

	return nil
}
