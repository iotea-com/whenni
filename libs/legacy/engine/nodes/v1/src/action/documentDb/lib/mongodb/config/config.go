package mongodbActionNodeConfig

import (
	"github.com/iotea-com/iotea/libs/engine/dependencies/models"
	"github.com/iotea-com/iotea/libs/engine/dependencies/things"
)

type MongoDbActionSubnodeConfig struct {
	MongoDb       things.MongoDbServer `json:"mongoDbServer::thing" validate:"required"`
	Database      string               `json:"database" validate:"required"`
	Collection    string               `json:"collection" validate:"required"`
	Method        string               `json:"queryMethod" validate:"required"`
	Filter        string               `json:"filter" validate:"required_if=Method Find,required_if=Method FindOne,required_if=Method FindOneAndUpdate,required_if=Method FindOneAndReplace,required_if=Method FindOneAndDelete,required_if=Method UpdateOne,required_if=Method UpdateMany,required_if=Method ReplaceOne,required_if=Method DeleteOne,required_if=Method DeleteMany"`
	Document      string               `json:"document" validate:"required_if=Method InsertOne,required_if=Method InsertMany,required_if=Method UpdateOne,required_if=Method UpdateMany,required_if=Method ReplaceOne,required_if=Method FindOneAndReplace,required_if=Method FindOneAndUpdate"`
	ModelInput    models.Model         `json:"input::model"`
	DefaultValues map[string]any       `json:"defaultValues"`
}
