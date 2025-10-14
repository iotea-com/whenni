package main

import (
	"context"
	"fmt"

	"github.com/goccy/go-json"

	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	mongodbActionNode "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/action/documentDb/lib/mongodb"
	mongodbActionNodeConfig "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/action/documentDb/lib/mongodb/config"
)

func (n *DocumentDbActionNode) Init(params node.InitParams) node.Error {
	n.logger = params.Obsv.Logger.With().Str("node", NodeName).Logger()

	// Attempt to create a MongoDB action subnode config
	mongoDbConfig := &mongodbActionNodeConfig.MongoDbActionSubnodeConfig{}
	err := json.Unmarshal(params.Config, mongoDbConfig)
	if err == nil && mongoDbConfig.MongoDb.Host != "" {
		n.logger.Debug().Msg("Recognized MongoDB configuration for document DB action node - using MongoDB action subnode")

		// New subnode
		n.subnode = mongodbActionNode.New()
		n.connectIoChannels()
		return n.subnode.Init(params)
	}

	return node.Error{
		Type:   node.ValidationError,
		Reason: fmt.Sprintf("unrecognized subnode config %s", string(params.Config)),
	}
}

func (n *DocumentDbActionNode) connectIoChannels() error {
	if len(n.subnode.GetInputChans(context.Background())) <= 0 {
		return fmt.Errorf("subnode did not have any input channels")
	}

	n.inputChannels[0] = n.subnode.GetInputChans(context.Background())[0]

	if len(n.subnode.GetOutputChans(context.Background())) <= 0 {
		return fmt.Errorf("subnode did not have any output channels")
	}

	n.outputChannels[0] = n.subnode.GetOutputChans(context.Background())[0]

	return nil
}
