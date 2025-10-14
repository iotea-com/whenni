package natsSourceNodeConfig

import (
	"github.com/iotea-com/iotea/libs/engine/dependencies/things"
)

type NatsSourceSubnodeConfig struct {
	ThingNatsClient things.NatsClient `json:"natsClient::thing" validate:"required"`
	Topic           string            `json:"topic" validate:"required"`
}
