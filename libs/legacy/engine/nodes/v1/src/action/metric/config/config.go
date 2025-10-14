package metricActionNodeConfig

import (
	"github.com/iotea-com/iotea/libs/engine/dependencies/models"
	"github.com/iotea-com/iotea/libs/engine/dependencies/things"
)

type MetricActionNodeConfig struct {
	ThingClickhouseDatabase things.ClickhouseDatabase `json:"clickhouseDatabase::thing" validate:"required"`
	Model                   models.Model              `json:"model" validate:"required"`
	Attribute               string                    `json:"attribute" validate:"required"`
	Metadata                string                    `json:"metadata" validate:"required"`
}
