package awsSESActionNodeConfig

import (
	"github.com/iotea-com/iotea/libs/engine/dependencies/models"
	"github.com/iotea-com/iotea/libs/engine/dependencies/things"
)

type AwsSESActionSubnodeConfig struct {
	ThingAwsSES     things.AwsSES  `json:"awsSES::thing" validate:"required"`
	Sender          string         `json:"sender" validate:"required"`
	Receiver        string         `json:"receiver" validate:"required"`
	Subject         string         `json:"subject" validate:"required"`
	ModelInput      *models.Model  `json:"input::model" validate:"required"`
	TemplateMessage string         `json:"templateMessage" validate:"required"`
	DefaultValues   map[string]any `json:"defaultValues" validate:"required"`
}
