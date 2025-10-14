package natsActionNode

import (
	"fmt"

	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
)

func (n *NatsActionSubnode) Exec(params node.ExecParams) node.Error {
	// Client library works with a registered callback. In this loop we simply
	// await for a context cancellation, or an error
Loop:
	for {
		select {
		case <-params.Ctx.Done():
			break Loop
		case ioData := <-n.inputChannels[0].Channel:
			invalid, err := n.handlePayload(ioData)
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

func (n *NatsActionSubnode) handlePayload(ioData node.IoData) (bool, error) {
	// Verify that there is a request body
	if ioData.Data == nil {
		return true, fmt.Errorf("no data was passed in input channel")
	}

	// Validate that the data is a []byte
	payloadData, ok := ioData.Data.([]byte)
	if !ok {
		return true, fmt.Errorf("expected data of type []byte, got %T", ioData.Data)
	}

	n.logger.Info().Ctx(ioData.Ctx).
		Str("data", string(payloadData)).
		Msgf("Publishing data to NATS topic %v", n.config.Topic)

	// Send the message
	err := n.conn.Publish(n.config.Topic, payloadData)
	if err != nil {
		return false, fmt.Errorf("error publishing data to NATS: %v", err)
	}

	// Flush to ensure the message is sent
	err = n.conn.Flush()
	if err != nil {
		return false, fmt.Errorf("error flushing the connection: %s", err)
	}

	// Check if there were any issues during publishing
	if err := n.conn.LastError(); err != nil {
		return false, fmt.Errorf("error after publish: %s", err)
	}

	n.logger.Info().Ctx(ioData.Ctx).Msgf("Successfully published data to NATS topic %v", n.config.Topic)

	return false, nil
}
