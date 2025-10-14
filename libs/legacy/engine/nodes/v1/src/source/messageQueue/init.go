package main

import (
	"context"
	"fmt"

	"github.com/goccy/go-json"

	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	kafkaSourceNode "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/source/messageQueue/lib/kafka"
	kafkaSourceNodeConfig "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/source/messageQueue/lib/kafka/config"
	natsSourceNode "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/source/messageQueue/lib/nats"
	natsSourceNodeConfig "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/source/messageQueue/lib/nats/config"
)

func (n *MessageQueueSourceNode) Init(params node.InitParams) node.Error {
	n.logger = params.Obsv.Logger.With().Str("node", NodeName).Logger()

	// Attempt to create a Kafka source subnode config
	kafkaConfig := &kafkaSourceNodeConfig.KafkaSourceSubnodeConfig{}
	err := json.Unmarshal(params.Config, kafkaConfig)
	if err == nil && kafkaConfig.ThingKafkaConsumer.Cluster != nil {
		n.logger.Debug().Msg("Recognized Kafka configuration for message queue source node - using Kafka source subnode")

		// new node
		n.subnode = kafkaSourceNode.New()
		n.connectIoChans()
		return n.subnode.Init(params)
	}

	// Attempt to create a NATS source subnode config
	natsConfig := &natsSourceNodeConfig.NatsSourceSubnodeConfig{}
	err = json.Unmarshal(params.Config, natsConfig)
	if err == nil && natsConfig.ThingNatsClient.Server != nil {
		n.logger.Debug().Msg("Recognized NATS configuration for message queue source node - using NATS source subnode")

		// new node
		n.subnode = natsSourceNode.New()
		n.connectIoChans()
		return n.subnode.Init(params)
	}

	return node.Error{
		Type:   node.ValidationError,
		Reason: fmt.Sprintf("unrecognized subnode config - %v: %s", err, string(params.Config)),
	}
}

func (n *MessageQueueSourceNode) connectIoChans() error {
	if len(n.subnode.GetInputChans(context.Background())) <= 0 {
		return fmt.Errorf("subnode did not have any input channels")
	}

	if len(n.subnode.GetOutputChans(context.Background())) <= 0 {
		return fmt.Errorf("subnode did not have any output channels")
	}

	n.inputChannels[0] = n.subnode.GetInputChans(context.Background())[0]
	n.outputChannels[0] = n.subnode.GetOutputChans(context.Background())[0]

	return nil
}
