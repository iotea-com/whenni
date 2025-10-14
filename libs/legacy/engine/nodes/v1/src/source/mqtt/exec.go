package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/iotea-com/iotea/libs/engine/logs"
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
)

func (n *MqttSourceNode) Exec(params node.ExecParams) node.Error {
	// Client library works with a registered callback. In this loop we simply
	// await for a context cancellation, or an error
Loop:
	for {
		select {
		case <-params.Ctx.Done():
			break Loop
		case err := <-n.errChan:
			return node.Error{
				Type:   node.FatalError,
				Reason: err.Error(),
			}

		}
	}
	return node.Error{
		Type: node.NoError,
	}
}

// mqttSubscriptionCallback will be called by the paho mqtt client library when a new message
// is received on the user subscribed topic. This function does not do any checks on the data,
// and simply forwards it to its output channel as a byte array
func (n *MqttSourceNode) mqttSubscriptionCallback(client mqtt.Client, msg mqtt.Message) {
	// Create an unique context for the request
	dataCtx := logs.NewDataCtx()

	// Notify the runtime that data was received
	n.notifyChannel <- node.Notification{
		Type:    node.NotifyInputDataReceived,
		DataCtx: dataCtx,
	}

	n.logger.Info().Ctx(dataCtx).Str("data", string(msg.Payload())).Msg("Message received")

	// Create the node data type
	outputData := node.IoData{
		Data: msg.Payload(),
		Type: n.outputChannels[0].Type[0], //TODO: fix to the right data type
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
