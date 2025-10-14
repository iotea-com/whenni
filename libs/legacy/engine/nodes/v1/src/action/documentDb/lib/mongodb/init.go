package mongodbActionNode

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/goccy/go-json"

	"github.com/agoda-com/opentelemetry-go/otelzerolog"
	"github.com/go-playground/validator/v10"
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	"github.com/iotea-com/iotea/libs/engine/template"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (n *MongoDbActionSubnode) Init(params node.InitParams) node.Error {
	startTime := time.Now()

	// Setup tracer, meter and logger
	if err := n.setupObservability(params.Obsv); err != nil {
		return node.Error{
			Type:   node.FatalError,
			Reason: fmt.Sprintf("error setting up observability: %v", err),
		}
	}

	// Set the notification channel
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
// model, and validates all the configuration fields
func (n *MongoDbActionSubnode) setConfig(config []byte) error {
	if err := json.Unmarshal(config, &n.config); err != nil {
		return fmt.Errorf("invalid node configuration JSON passed: %v", err)
	}

	// Validate config
	v := validator.New()

	if err := v.Struct(&n.config); err != nil {
		return err
	}

	// Validate model input schema
	if err := n.config.ModelInput.Validate(); err != nil {
		return fmt.Errorf("invalid model input schema: %v", err)
	}

	// Create & verify templates
	n.filterTemplate = template.Template{
		Value:         n.config.Filter,
		ModelInput:    &n.config.ModelInput,
		DefaultValues: n.config.DefaultValues,
	}

	switch n.config.Method {
	case string(MongoDbQueryMethodFind),
		string(MongoDbQueryMethodFindOne),
		string(MongoDbQueryMethodFindOneAndDelete),
		string(MongoDbQueryMethodFindOneAndReplace),
		string(MongoDbQueryMethodFindOneAndUpdate),
		string(MongoDbQueryMethodUpdateOne),
		string(MongoDbQueryMethodUpdateMany),
		string(MongoDbQueryMethodReplaceOne),
		string(MongoDbQueryMethodDeleteOne),
		string(MongoDbQueryMethodDeleteMany):
		if err := n.filterTemplate.Verify(); err != nil {
			return fmt.Errorf("error verifying filter template: %s", err)
		}
	}

	n.documentTemplate = template.Template{
		Value:         n.config.Document,
		ModelInput:    &n.config.ModelInput,
		DefaultValues: n.config.DefaultValues,
	}

	switch n.config.Method {
	case string(MongoDbQueryMethodInsertOne),
		string(MongoDbQueryMethodInsertMany),
		string(MongoDbQueryMethodUpdateOne),
		string(MongoDbQueryMethodUpdateMany),
		string(MongoDbQueryMethodReplaceOne),
		string(MongoDbQueryMethodFindOneAndReplace),
		string(MongoDbQueryMethodFindOneAndUpdate):
		if err := n.documentTemplate.Verify(); err != nil {
			return fmt.Errorf("error verifying document template: %s", err)
		}
	}

	// Connect to MongoDB
	databaseUrl := (func() string {
		if n.config.MongoDb.Protocol == "mongodb" {
			return fmt.Sprintf("%s://%s:%s@%s:%d/", n.config.MongoDb.Protocol, n.config.MongoDb.Username, n.config.MongoDb.Password, n.config.MongoDb.Host, n.config.MongoDb.Port)
		}

		return fmt.Sprintf("%s://%s:%s@%s/", n.config.MongoDb.Protocol, n.config.MongoDb.Username, n.config.MongoDb.Password, n.config.MongoDb.Host)
	})()
	clientOptions := options.Client().ApplyURI(databaseUrl)

	dbCtx := context.Background()
	defer dbCtx.Done()

	client, err := mongo.Connect(dbCtx, clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to a MongoDB server at %s:%d: %s", n.config.MongoDb.Host, n.config.MongoDb.Port, err)
	}

	err = client.Ping(dbCtx, nil)
	if err != nil {
		return fmt.Errorf("failed to ping the MongoDB server at %s:%d: %s", n.config.MongoDb.Host, n.config.MongoDb.Port, err)
	}

	// Check database existence
	databaseNames, err := client.ListDatabaseNames(dbCtx, bson.M{})
	if err != nil {
		return fmt.Errorf("failed to get database names")
	}

	if databaseNames == nil || !slices.Contains(databaseNames, n.config.Database) {
		return fmt.Errorf("database '%s' was not found", n.config.Database)
	}

	// TODO: Check collection existence
	// database := n.mongoClient.Database(n.config.Database)
	// collectionNames, err := database.ListCollectionNames(dbCtx, bson.M{})
	// if err != nil {
	// 	n.logger.Error().Str("node", NodeName).Ctx(logCtx).
	// 		Msgf("failed to get collection names in database '%s' - a collection will be created if it does not exist", n.config.Database)
	// }

	// if collectionNames == nil || !slices.Contains(collectionNames, n.config.Collection) {
	// 	n.logger.Error().Str("node", NodeName).Ctx(logCtx).
	// 		Msgf("collection '%s' was not found in database '%s' - a collection will be created if it does not exist", n.config.Collection, n.config.Database)
	// }

	n.mongoClient = client

	return nil
}

// setupObservability will create a tracer, a meter and a logger to be used
// internally by the node, based on the providers passed at initialization
func (n *MongoDbActionSubnode) setupObservability(obsv node.Observability) error {
	n.tracer = obsv.TracerProvider.Tracer("MongoDbActionNodeTracer")
	n.meter = obsv.MeterProvider.Meter("MongoDbActionNodeMeter")

	hook := otelzerolog.NewHook(obsv.LoggerProvider)
	n.logger = obsv.Logger.With().Str("node", NodeName).Logger().Hook(hook)

	return nil
}
