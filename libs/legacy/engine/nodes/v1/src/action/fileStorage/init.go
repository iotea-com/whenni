package main

import (
	"context"
	"fmt"

	"github.com/goccy/go-json"

	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	awsS3Node "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/action/fileStorage/lib/awsS3"
	s3ActionNodeConfig "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/action/fileStorage/lib/awsS3/config"
	minioNode "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/action/fileStorage/lib/minio"
	minioActionNodeConfig "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/action/fileStorage/lib/minio/config"
)

func (n *FileStorageActionNode) Init(params node.InitParams) node.Error {
	n.logger = params.Obsv.Logger.With().Str("node", NodeName).Logger()

	// Attempt to create a Minio action subnode config
	minioConfig := &minioActionNodeConfig.MinioActionSubnodeConfig{}
	err := json.Unmarshal(params.Config, minioConfig)
	if err == nil && minioConfig.ThingMinioBucket.MinioSecretAccessKey != "" {
		n.logger.Debug().Msg("Recognized Minio configuration for file storage action node - using Minio action subnode")

		// new node
		n.subnode = minioNode.New()
		n.connectIoChannels()
		return n.subnode.Init(params)
	}
	// Attempt to create an AWS S3 action subnode config
	s3Config := &s3ActionNodeConfig.S3ActionSubnodeConfig{}
	err = json.Unmarshal(params.Config, s3Config)
	if err == nil && s3Config.ThingAwsS3Bucket.AwsSecretAccessKey != "" {
		n.logger.Debug().Msg("Recognized AWS S3 configuration for file storage action node - using S3 action subnode")

		// new node
		n.subnode = awsS3Node.New()
		n.connectIoChannels()
		return n.subnode.Init(params)
	}

	return node.Error{
		Type:   node.ValidationError,
		Reason: fmt.Sprintf("unrecognized subnode config - %s: %s", err, params.Config),
	}
}

func (n *FileStorageActionNode) connectIoChannels() error {
	if len(n.subnode.GetInputChans(context.Background())) <= 0 {
		return fmt.Errorf("subnode did not have any input channels")
	}
	n.inputChannels[0] = n.subnode.GetInputChans(context.Background())[0]
	return nil
}
