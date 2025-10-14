package main

import (
	"github.com/goccy/go-json"

	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
)

func (n *LogActionNode) Exec(params node.ExecParams) node.Error {
Loop:
	for {
		select {
		case <-params.Ctx.Done():
			break Loop
		case ioData := <-n.inputChannels[0].Channel:
			switch ioData.Type {
			case node.BytesDataType:
				var j map[string]any
				err := json.Unmarshal(ioData.Data.([]byte), &j)
				if err != nil {
					n.logger.Info().Ctx(ioData.Ctx).Msgf("Data (type %d): %s", ioData.Type, ioData.Data)
				} else {
					prettyJsonString, _ := json.MarshalIndent(j, "", "  ")

					n.logger.Info().Ctx(ioData.Ctx).Msgf("Data (type %d): %s", ioData.Type, prettyJsonString)
				}
			case node.MapDataType:
				prettyJsonString, err := json.MarshalIndent(ioData.Data, "", "  ")
				if err != nil {
					n.logger.Info().Ctx(ioData.Ctx).Msgf("Could not marshal input data: %s", err)
				} else {
					n.logger.Info().Ctx(ioData.Ctx).Msgf("Data (type %d): %s", ioData.Type, prettyJsonString)
				}
			case node.StringDataType:
				n.logger.Info().Ctx(ioData.Ctx).Msgf("Data (type %d): %s", ioData.Type, ioData.Data)
			case node.SignalDataType:
				n.logger.Info().Ctx(ioData.Ctx).Msgf("Signal received")
			default:
				n.logger.Info().Ctx(ioData.Ctx).Msgf("Data (type %d): %v", ioData.Type, ioData.Data)
			}

			// Safely notify the runtime that data was processed
			n.notifyChannel <- node.Notification{
				Type:    node.NotifyDataProcessed,
				DataCtx: ioData.Ctx,
			}
		}
	}
	return node.Error{
		Type: node.NoError,
	}
}
