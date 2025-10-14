package main

import (
	"context"
	"fmt"

	"github.com/goccy/go-json"

	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	awsSESActionNode "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/action/notification/lib/awsSES"
	awsSESActionNodeConfig "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/action/notification/lib/awsSES/config"
	awsSNSActionNode "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/action/notification/lib/awsSNS"
	awsSNSActionNodeConfig "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/action/notification/lib/awsSNS/config"
	sendgridActionNode "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/action/notification/lib/sendgrid"
	sendgridActionNodeConfig "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/action/notification/lib/sendgrid/config"
)

func (n *NotificationActionNode) Init(params node.InitParams) node.Error {
	n.logger = params.Obsv.Logger.With().Str("node", NodeName).Logger()

	// Attempt to create a AWS SNS action subnode config
	awsSNSConfig := &awsSNSActionNodeConfig.AwsSNSActionSubnodeConfig{}
	err := json.Unmarshal(params.Config, awsSNSConfig)
	if err == nil && awsSNSConfig.ThingAwsSNS.AwsAccessKeyId != "" && awsSNSConfig.TopicARN != "" {
		n.logger.Debug().Msg("Recognized AWS SNS configuration for notification action node - using AWS SNS action subnode")

		// new node
		n.subnode = awsSNSActionNode.New()
		n.connectIoChannels()
		return n.subnode.Init(params)
	}

	// Attempt to create a AWS SES action subnode config
	awsSESConfig := &awsSESActionNodeConfig.AwsSESActionSubnodeConfig{}
	err = json.Unmarshal(params.Config, awsSESConfig)
	if err == nil && awsSESConfig.ThingAwsSES.AwsAccessKeyId != "" && awsSESConfig.Sender != "" {
		n.logger.Debug().Msg("Recognized AWS SES configuration for notification action node - using AWS SES action subnode")

		// new node
		n.subnode = awsSESActionNode.New()
		n.connectIoChannels()
		return n.subnode.Init(params)
	}

	// Attempt to create a Sendgrid action subnode config
	sendgridConfig := &sendgridActionNodeConfig.SendgridActionSubnodeConfig{}
	err = json.Unmarshal(params.Config, sendgridConfig)
	if err == nil && sendgridConfig.SendgridClient.ApiKey != "" {
		n.logger.Debug().Msg("Recognized Sendgrid configuration for notification action node - using Sendgrid action subnode")

		// new node
		n.subnode = sendgridActionNode.New()
		n.connectIoChannels()
		return n.subnode.Init(params)
	}

	return node.Error{
		Type:   node.ValidationError,
		Reason: fmt.Sprintf("unrecognized subnode config - %s: %s", err, params.Config),
	}
}

func (n *NotificationActionNode) connectIoChannels() error {
	if len(n.subnode.GetInputChans(context.Background())) <= 0 {
		return fmt.Errorf("subnode did not have any input channels")
	}

	n.inputChannels[0] = n.subnode.GetInputChans(context.Background())[0]

	return nil
}
