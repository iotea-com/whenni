package minioActionNodeConfig

import (
	"github.com/iotea-com/iotea/libs/engine/dependencies/models"
	"github.com/iotea-com/iotea/libs/engine/dependencies/things"
)

type MinioActionSubnodeConfig struct {
	ThingMinioBucket things.MinioBucket `json:"minioBucket::thing" validate:"required"`
	Key              string             `json:"key" validate:"required"`
	ModelInput       models.Model       `json:"input::model"`
	DefaultValues    map[string]any     `json:"defaultValues"`
}
