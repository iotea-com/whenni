package main

import (
	"context"
	"fmt"

	"github.com/goccy/go-json"

	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	influxdbNode "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/action/timeSeriesDb/lib/influxdb"
	influxdbActionNodeConfig "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/action/timeSeriesDb/lib/influxdb/config"
)

func (n *TimeSeriesDbActionNode) Init(params node.InitParams) node.Error {
	n.logger = params.Obsv.Logger.With().Str("node", NodeName).Logger()

	// Attempt to create a Kafka action subnode config
	influxdbConfig := &influxdbActionNodeConfig.InfluxDBActionNodeConfig{}
	err := json.Unmarshal(params.Config, influxdbConfig)
	if err == nil {
		n.logger.Debug().Msg("Recognized InfluxDB configuration for time series database action node - using InfluxDB action subnode")

		// new node
		n.subnode = influxdbNode.New()
		n.connectIoChannels()
		return n.subnode.Init(params)
	}

	return node.Error{
		Type:   node.ValidationError,
		Reason: fmt.Sprintf("unrecognized subnode config - %s: %s", err, params.Config),
	}
}

func (n *TimeSeriesDbActionNode) connectIoChannels() error {
	if len(n.subnode.GetInputChans(context.Background())) <= 0 {
		return fmt.Errorf("subnode did not have any input channels")
	}

	n.inputChannels[0] = n.subnode.GetInputChans(context.Background())[0]

	return nil
}
