package main

import (
	"github.com/iotea-com/iotea/libs/engine/logs"
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
)

func (n *TimerSourceNode) Exec(params node.ExecParams) node.Error {
	// Start the cron scheduler
	n.cronScheduler.Start()

Loop:
	for {
		select {
		case <-params.Ctx.Done():
			break Loop
		case <-n.triggerChannel:
			// Create an unique context for the request
			dataCtx := logs.NewDataCtx()

			// Notify the runtime that data was received
			n.notifyChannel <- node.Notification{
				Type:    node.NotifyInputDataReceived,
				DataCtx: dataCtx,
			}

			n.logger.Info().Ctx(dataCtx).Msg("Timer triggered")

			// Send a signal to the output channel
			outputData := node.IoData{
				Data: nil,
				Type: node.SignalDataType,
				Ctx:  dataCtx,
			}

			// Send data to the next node or the trigger
			n.outputChannels[0].Channel <- outputData

			// Notify the runtime that the data was processed
			n.notifyChannel <- node.Notification{
				Type:    node.NotifyDataProcessed,
				DataCtx: dataCtx,
			}
		}
	}

	return node.Error{
		Type: node.NoError,
	}
}
