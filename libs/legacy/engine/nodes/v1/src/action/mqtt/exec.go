package main

import (
	"fmt"

	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
)

func (n *MqttActionNode) Exec(params node.ExecParams) node.Error {
Loop:
	for {
		select {
		case <-params.Ctx.Done():
			break Loop
		case ioData := <-n.inputChannels[0].Channel:
			invalid, err := n.publishMessage(ioData)
			if err != nil {
				if invalid {
					// Safely notify the runtime that data was processed
					n.notifyChannel <- node.Notification{
						Type:    node.NotifyDataInvalid,
						DataCtx: ioData.Ctx,
						Reason:  err.Error(),
					}
				} else {
					return node.Error{
						Type:   node.FatalError,
						Reason: err.Error(),
					}
				}
			} else {
				// Safely notify the runtime that data was processed
				n.notifyChannel <- node.Notification{
					Type:    node.NotifyDataProcessed,
					DataCtx: ioData.Ctx,
				}
			}
		}
	}

	return node.Error{
		Type: node.NoError,
	}
}

// publishMessage publishes a message to the MQTT broker.
func (n *MqttActionNode) publishMessage(ioData node.IoData) (bool, error) {
	payloadData, ok := ioData.Data.([]byte)
	if !ok {
		return true, fmt.Errorf("expected data of type []byte, got %T", ioData.Data) // Validation error
	}

	n.logger.Info().Ctx(ioData.Ctx).
		Str("data", string(payloadData)).
		Msgf("Publishing data to '%v'", n.config.Topic)

	// Publish data to MQTT topic
	token := n.client.Publish(n.config.Topic, byte(n.config.QoS), false, payloadData)
	token.Wait()
	if token.Error() != nil {
		return false, fmt.Errorf("error publishing data to MQTT: %v", token.Error()) // Fatal error
	} else {
		n.logger.Info().Ctx(ioData.Ctx).Msgf("Data published to '%v' successfully", n.config.Topic)

		return false, nil // No error
	}
}
