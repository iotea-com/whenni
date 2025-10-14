package influxdbNode

import (
	"context"
	"fmt"
	"time"

	"github.com/goccy/go-json"

	"github.com/agoda-com/opentelemetry-go/otelzerolog"
	"github.com/go-playground/validator/v10"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/domain"
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
)

func (n *InfluxDBActionNode) Init(params node.InitParams) node.Error {
	startTime := time.Now()

	// Setup tracer, meter and logger
	if err := n.setupObservability(params.Obsv); err != nil {
		return node.Error{
			Type:   node.FatalError,
			Reason: err.Error(),
		}
	}

	// Set the notifications channel
	n.notifyChannel = params.NotifyChannel

	// Setup configuration fields
	if err := n.setConfig(params.Config); err != nil {
		return node.Error{
			Type:   node.ValidationError,
			Reason: fmt.Sprintf("invalid configuration, %v", err),
		}
	}

	// Initialize the Database
	if err := n.initDatabase(); err != nil {
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

// setupObservability will create a tracer, a meter and a logger to be used
// internally by the node, based on the providers passed at initialization
func (n *InfluxDBActionNode) setupObservability(obsv node.Observability) error {
	n.tracer = obsv.TracerProvider.Tracer("DataMappingProcessingNodeTracer")
	n.meter = obsv.MeterProvider.Meter("DataMappingProcessingNodeMeter")

	hook := otelzerolog.NewHook(obsv.LoggerProvider)
	n.logger = obsv.Logger.With().Str("node", NodeName).Logger().Hook(hook)

	return nil
}

// setConfig unmarshalls the passed config byte array into the node's configuration
// struct, and validates all the configuration fields
func (n *InfluxDBActionNode) setConfig(config []byte) error {
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

// initialize checks for the existence of a bucket and creates it if it does not exist
func (n *InfluxDBActionNode) initDatabase() error {
	databaseUrl := fmt.Sprintf("%s://%s:%d", n.config.ThingInfluxDB.Protocol, n.config.ThingInfluxDB.Host, n.config.ThingInfluxDB.Port)
	n.client = influxdb2.NewClient(databaseUrl, n.config.ThingInfluxDB.Token)

	ctx := context.Background()
	bucketsAPI := n.client.BucketsAPI()
	orgsAPI := n.client.OrganizationsAPI()

	// Get the OrgID from the OrgName
	org, err := orgsAPI.FindOrganizationByName(ctx, n.config.ThingInfluxDB.OrgName)
	if err != nil {
		return fmt.Errorf("error finding organization %s: %v", n.config.ThingInfluxDB.OrgName, err)
	}
	if org == nil {
		return fmt.Errorf("organization %s not found", n.config.ThingInfluxDB.OrgName)
	}

	// Check if the bucket exists
	bucket, err := bucketsAPI.FindBucketByName(ctx, n.config.Bucket)
	if err != nil && err.Error() != fmt.Sprintf("bucket '%s' not found", n.config.Bucket) {
		return fmt.Errorf("error checking bucket: %v", err)
	}
	if bucket == nil { // Bucket not found, let's create it
		retentionType := domain.RetentionRuleType("forever")
		newBucket := &domain.Bucket{
			OrgID: org.Id,
			Name:  n.config.Bucket,
			RetentionRules: []domain.RetentionRule{{
				Type:         &retentionType,
				EverySeconds: 0, // No retention period, meaning data is kept forever
			}},
		}
		if _, err := bucketsAPI.CreateBucket(ctx, newBucket); err != nil {
			return fmt.Errorf("failed to create bucket: %v", err)
		}

		n.logger.Debug().Msgf("Created '%v' bucket in database instance", n.config.Bucket)
	} else {
		n.logger.Debug().Msgf("Node using '%v' bucket, already present in database instance", n.config.Bucket)
	}

	// Instantiate the write API for this organization and bucket
	n.writeAPI = n.client.WriteAPIBlocking(*org.Id, n.config.Bucket)

	return nil
}
