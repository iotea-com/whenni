package mongodbActionNode

import (
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	mongodbActionNodeConfig "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/action/documentDb/lib/mongodb/config"
	"github.com/iotea-com/iotea/libs/engine/template"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

const (
	NodeName = "MongoDB Client"
)

type MongoDbQueryMethod string

const (
	MongoDbQueryMethodFind              MongoDbQueryMethod = "Find"
	MongoDbQueryMethodFindOne           MongoDbQueryMethod = "FindOne"
	MongoDbQueryMethodFindOneAndUpdate  MongoDbQueryMethod = "FindOneAndUpdate"
	MongoDbQueryMethodFindOneAndReplace MongoDbQueryMethod = "FindOneAndReplace"
	MongoDbQueryMethodFindOneAndDelete  MongoDbQueryMethod = "FindOneAndDelete"
	MongoDbQueryMethodInsertOne         MongoDbQueryMethod = "InsertOne"
	MongoDbQueryMethodInsertMany        MongoDbQueryMethod = "InsertMany"
	MongoDbQueryMethodUpdateOne         MongoDbQueryMethod = "UpdateOne"
	MongoDbQueryMethodUpdateMany        MongoDbQueryMethod = "UpdateMany"
	MongoDbQueryMethodReplaceOne        MongoDbQueryMethod = "ReplaceOne"
	MongoDbQueryMethodDeleteOne         MongoDbQueryMethod = "DeleteOne"
	MongoDbQueryMethodDeleteMany        MongoDbQueryMethod = "DeleteMany"
)

type MongoDbActionSubnode struct {
	config         mongodbActionNodeConfig.MongoDbActionSubnodeConfig
	inputChannels  []node.IoChannel
	outputChannels []node.IoChannel
	notifyChannel  chan node.Notification

	filterTemplate   template.Template
	documentTemplate template.Template

	mongoClient *mongo.Client

	// Observability
	tracer trace.Tracer
	meter  metric.Meter
	logger zerolog.Logger
}

func New() node.Interface {
	return &MongoDbActionSubnode{
		inputChannels: []node.IoChannel{
			{
				Id:      "documentDbInput",
				Channel: make(chan node.IoData, 100),
				Type:    []node.DataType{node.BytesDataType},
			},
		},
		outputChannels: []node.IoChannel{
			{
				Id:      "documentDbOutput",
				Channel: make(chan node.IoData, 100),
				Type:    []node.DataType{node.BytesDataType},
			},
		},
	}
}
