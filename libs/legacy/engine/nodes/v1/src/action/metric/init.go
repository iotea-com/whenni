package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/goccy/go-json"

	"github.com/agoda-com/opentelemetry-go/otelzerolog"
	"github.com/go-playground/validator/v10"
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
)

func (n *MetricActionNode) Init(params node.InitParams) node.Error {
	startTime := time.Now()

	// Setup tracer, meter and logger
	if err := n.setupObservability(params.Obsv); err != nil {
		return node.Error{
			Type:   node.FatalError,
			Reason: err.Error(),
		}
	}

	// A notification channel is used to tell the runtime about the
	// flow of data through the node.
	n.notifyChannel = params.NotifyChannel

	if err := n.setConfig(params.Config); err != nil {
		return node.Error{
			Type:   node.ValidationError,
			Reason: fmt.Sprintf("invalid configuration, %v", err),
		}
	}

	if err := n.connectToClickhouse(); err != nil {
		return node.Error{
			Type:   node.FatalError,
			Reason: fmt.Sprintf("could not connect to ClickHouse: %s", err),
		}
	}

	n.logger.Debug().Msgf("%s Node Initialized OK in %vms", NodeName, time.Since(startTime).Milliseconds())

	return node.Error{
		Type: node.NoError,
	}
}

// setConfig unmarshals the passed config byte array into the node's configuration
// struct, and validates all the configuration fields
func (n *MetricActionNode) setConfig(config []byte) error {
	if err := json.Unmarshal(config, &n.config); err != nil {
		return fmt.Errorf("invalid node configuration JSON passed: %v", err)
	}

	// Validate config
	v := validator.New()

	if err := v.Struct(&n.config); err != nil {
		return err
	}

	// Validate model configuration
	if err := n.checkConfig(); err != nil {
		return err
	}

	return nil
}

// checkConfig validates that the model configuration matches our requirements:
// - Must have a numeric field for the value
// - Must have a field for metadata
func (n *MetricActionNode) checkConfig() error {
	// Check if the value field exists and is numeric
	valueFieldId := n.config.Attribute
	if valueFieldId == "" {
		return fmt.Errorf("value field must be specified")
	}

	// Find the value field in the model attributes
	valueAttr, exists := n.config.Model.Attributes[valueFieldId]
	if !exists {
		return fmt.Errorf("value field ID '%s' not found in model attributes", valueFieldId)
	}
	if valueAttr.Type != "number" {
		return fmt.Errorf("value field '%s' must be of type 'number', got '%s'", valueAttr.Key, valueAttr.Type)
	}

	// Check if the metadata field exists
	metadataFieldId := n.config.Metadata
	if metadataFieldId == "" {
		return fmt.Errorf("metadata field must be specified")
	}

	// Find the metadata field in the model attributes
	metadataAttr, exists := n.config.Model.Attributes[metadataFieldId]
	if !exists {
		return fmt.Errorf("metadata field ID '%s' not found in model attributes", metadataFieldId)
	}
	if !isValidMetadataType(metadataAttr.Type) {
		return fmt.Errorf("metadata field '%s' must be of type 'string', 'number', or 'boolean', got '%s'", metadataAttr.Key, metadataAttr.Type)
	}

	return nil
}

// isValidMetadataType checks if the given type is valid for metadata
func isValidMetadataType(attrType string) bool {
	validTypes := map[string]bool{
		"string":  true,
		"number":  true,
		"boolean": true,
	}
	return validTypes[attrType]
}

// setupObservability will create a tracer, a meter and a logger to be used
// internally by the node, based on the providers passed at initialization
func (n *MetricActionNode) setupObservability(obsv node.Observability) error {
	n.tracer = obsv.TracerProvider.Tracer("MetricActionNodeTracer")
	n.meter = obsv.MeterProvider.Meter("MetricActionNodeMeter")

	hook := otelzerolog.NewHook(obsv.LoggerProvider)
	n.logger = obsv.Logger.With().Str("node", NodeName).Logger().Hook(hook)

	return nil
}

func (n *MetricActionNode) connectToClickhouse() error {
	// Connect to ClickHouse
	addr := fmt.Sprintf("%s:%d", n.config.ThingClickhouseDatabase.Host, n.config.ThingClickhouseDatabase.Port)
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr:     []string{addr},
		Protocol: clickhouse.Native,
		Auth: clickhouse.Auth{
			Database: n.config.ThingClickhouseDatabase.Database,
			Username: n.config.ThingClickhouseDatabase.Username,
			Password: n.config.ThingClickhouseDatabase.Password,
		},
		MaxOpenConns: 10,
		MaxIdleConns: 10,
	})
	if err != nil {
		return fmt.Errorf("could not connect to ClickHouse: %s", err)
	}

	// Give it some time to connect
	pingTimeout := 3 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
	defer cancel()

	// Ping the server
	err = conn.Ping(ctx)
	if err != nil {
		return fmt.Errorf("failed to ping the ClickHouse server at %s:%d: %s", n.config.ThingClickhouseDatabase.Host, n.config.ThingClickhouseDatabase.Port, err)
	}

	// Check that the database exists
	query := fmt.Sprintf("SELECT name FROM system.databases WHERE name = '%s'", n.config.ThingClickhouseDatabase.Database)
	var result string
	err = conn.QueryRow(ctx, query).Scan(&result)
	if err != nil {
		return fmt.Errorf("failed to check if database exists: %s", err)
	}

	// Check that the table exists
	query = fmt.Sprintf("SELECT name FROM system.tables WHERE name = '%s'", CustomMetricsTable)
	err = conn.QueryRow(ctx, query).Scan(&result)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("required table '%s' does not exist in database '%s'", CustomMetricsTable, n.config.ThingClickhouseDatabase.Database)
		}
		return fmt.Errorf("failed to check if table %s exists: %s", CustomMetricsTable, err)
	}

	// Great! now that we are connected, let's set the global client
	n.clickhouseClient = &conn

	n.logger.Info().Msgf("Connected to ClickHouse database %s", n.config.ThingClickhouseDatabase.Database)

	return nil
}
