package influxdbActionNodeConfig

import (
	"github.com/iotea-com/iotea/libs/engine/dependencies/things"
)

type InfluxDBActionNodeConfig struct {
	ThingInfluxDB things.InfluxDbDatabase `json:"influxdbDatabase::thing" validate:"required"`
	Bucket        string                  `json:"bucket" validate:"required"`
}
