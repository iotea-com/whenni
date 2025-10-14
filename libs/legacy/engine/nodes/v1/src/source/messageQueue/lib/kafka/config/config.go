package kafkaSourceNodeConfig

import (
	"github.com/iotea-com/iotea/libs/engine/dependencies/things"
)

type KafkaSourceSubnodeConfig struct {
	ThingKafkaConsumer things.KafkaConsumer `json:"kafkaConsumer::thing" validate:"required"`
	Topic              string               `json:"topic" validate:"required"`
}
