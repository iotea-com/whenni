package sendgridActionNodeConfig

import (
	"github.com/iotea-com/iotea/libs/engine/dependencies/models"
	"github.com/iotea-com/iotea/libs/engine/dependencies/things"
)

type SendgridActionSubnodeConfig struct {
	SendgridClient  things.SendgridClient `json:"sendgridClient::thing" validate:"required"`
	Sender          string                `json:"sender" validate:"required"`
	SendAs          *string               `json:"sendAs"`
	Receiver        string                `json:"receiver" validate:"required"`
	Subject         string                `json:"subject" validate:"required"`
	ModelInput      *models.Model         `json:"input::model"`
	TemplateMessage string                `json:"templateMessage" validate:"required"`
	DefaultValues   map[string]any        `json:"defaultValues"`
}
