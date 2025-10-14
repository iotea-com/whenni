package awsSNSActionNodeConfig

import (
	"github.com/iotea-com/iotea/libs/engine/dependencies/models"
	"github.com/iotea-com/iotea/libs/engine/dependencies/things"
)

type AwsSNSActionSubnodeConfig struct {
	ThingAwsSNS     things.AwsSNS  `json:"awsSNS::thing" validate:"required"`
	TopicARN        string         `json:"topicARN" validate:"required"`
	ModelInput      *models.Model  `json:"input::model" validate:"required"`
	TemplateMessage string         `json:"templateMessage" validate:"required"`
	DefaultValues   map[string]any `json:"defaultValues" validate:"required"`
}
