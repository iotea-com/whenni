package kafkaActionNodeConfig

import (
	"github.com/iotea-com/iotea/libs/engine/dependencies/things"
)

type KafkaActionSubnodeConfig struct {
	ThingKafkaProducer things.KafkaProducer `json:"kafkaProducer::thing" validate:"required"`
	Topic              string               `json:"topic" validate:"required"`
}
