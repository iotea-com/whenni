package awsS3ActionNodeConfig

import (
	"github.com/iotea-com/iotea/libs/engine/dependencies/models"
	"github.com/iotea-com/iotea/libs/engine/dependencies/things"
)

type S3ActionSubnodeConfig struct {
	ThingAwsS3Bucket things.S3Bucket `json:"awsS3Bucket::thing" validate:"required"`
	Key              string          `json:"key" validate:"required"`
	ModelInput       models.Model    `json:"input::model"`
	DefaultValues    map[string]any  `json:"defaultValues"`
}
