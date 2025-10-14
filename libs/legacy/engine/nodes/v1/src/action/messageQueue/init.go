package main

import (
	"context"
	"fmt"

	"github.com/goccy/go-json"

	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	kafkaNode "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/action/messageQueue/lib/kafka"
	kafkaActionNodeConfig "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/action/messageQueue/lib/kafka/config"
	natsNode "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/action/messageQueue/lib/nats"
	natsActionNodeConfig "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/action/messageQueue/lib/nats/config"
)

func (n *MessageQueueActionNode) Init(params node.InitParams) node.Error {
	n.logger = params.Obsv.Logger.With().Str("node", NodeName).Logger()

	// Attempt to create a Kafka action subnode config
	kafkaConfig := &kafkaActionNodeConfig.KafkaActionSubnodeConfig{}
	err := json.Unmarshal(params.Config, kafkaConfig)
	if err == nil && kafkaConfig.ThingKafkaProducer.Cluster != nil {
		n.logger.Debug().Msg("Recognized Kafka configuration for message queue action node - using Kafka action subnode")

		// new node
		n.subnode = kafkaNode.New()
		n.connectIoChannels()
		return n.subnode.Init(params)
	}

	// Attempt to create a NATS action subnode config
	natsConfig := &natsActionNodeConfig.NatsActionSubnodeConfig{}
	err = json.Unmarshal(params.Config, natsConfig)
	if err == nil && natsConfig.ThingNatsClient.Server != nil {
		n.logger.Debug().Msg("Recognized NATS configuration for message queue action node - using NATS action subnode")

		// new node
		n.subnode = natsNode.New()
		n.connectIoChannels()
		return n.subnode.Init(params)
	}

	return node.Error{
		Type:   node.ValidationError,
		Reason: fmt.Sprintf("unrecognized subnode config - %s: %s", err, params.Config),
	}
}

func (n *MessageQueueActionNode) connectIoChannels() error {
	if len(n.subnode.GetInputChans(context.Background())) <= 0 {
		return fmt.Errorf("subnode did not have any input channels")
	}

	n.inputChannels[0] = n.subnode.GetInputChans(context.Background())[0]

	return nil
}
