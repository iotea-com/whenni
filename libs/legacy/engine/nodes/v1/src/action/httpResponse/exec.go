package main

import (
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"

	"fmt"

	httpSourceNodeConfig "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/source/http/config"
)

func (n *HttpResponseActionNode) Exec(params node.ExecParams) node.Error {
Loop:
	for {
		select {
		case <-params.Ctx.Done(): // Rules engine exit signal
			break Loop
		case ioData := <-n.inputChannels[0].Channel: // Wait for input data
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

func (n *HttpResponseActionNode) handlePayload(ioData node.IoData) (bool, error) {
	if ioData.Data == nil {
		return true, fmt.Errorf("input data is nil")
	}

	// Safely convert data to bytes
	data, ok := ioData.Data.([]byte)
	if !ok {
		return true, fmt.Errorf("input data is not in expected format ([]byte)")
	}

	response := httpSourceNodeConfig.HttpResponseInfo{
		ResponseCode: n.config.ResponseCode,
		Headers:      n.config.Headers,
		ResponseBody: data,
	}

	outputData := node.IoData{
		Data: response,
		Type: node.MapDataType,
		Ctx:  ioData.Ctx,
	}

	select {
	case n.outputChannels[0].Channel <- outputData:
		return false, nil
	default:
		return false, fmt.Errorf("output channel is full")
	}
}
