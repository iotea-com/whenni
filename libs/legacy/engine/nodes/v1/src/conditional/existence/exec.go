package main

import (
	"fmt"

	"github.com/goccy/go-json"

	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
)

func (n *ExistenceConditionalNode) Exec(params node.ExecParams) node.Error {
Loop:
	for {
		select {
		case <-params.Ctx.Done():
			break Loop
		case ioData := <-n.inputChannels[0].Channel:
			output, conditionPass, invalid, err := n.handle(ioData)
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
				output.Ctx = params.Ctx
				if conditionPass {
					n.outputChannels[1].Channel <- *output
				} else {
					n.outputChannels[0].Channel <- *output
				}

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

// handle processes the input data and returns the output data, whether the condition passed,
// whether the data is invalid, and an error if any.
func (n *ExistenceConditionalNode) handle(ioData node.IoData) (*node.IoData, bool, bool, error) {
	// Ensure the input type is correct
	dataBytes, ok := ioData.Data.([]byte)
	if !ok {
		return nil, false, true, fmt.Errorf("input data has wrong type, expected []byte, got %T", ioData.Data)
	}

	n.logger.Info().Ctx(ioData.Ctx).Str("data", string(dataBytes)).Msg("Input data")

	// Initialize the input map
	input := make(map[string]json.RawMessage)

	// Unpack the input into a JSON struct
	err := json.Unmarshal(dataBytes, &input)
	if err != nil {
		return nil, false, true, fmt.Errorf("failed to unmarshal byte data: %s", err)
	}

	existenceResults := []bool{}
	// Check the existence of each attributeId
	for _, attributeId := range n.config.AttributeIds {
		if failureReason := n.config.Model.CheckExistence(input, attributeId); failureReason == nil {
			n.logger.Info().Ctx(ioData.Ctx).Msgf("Field '%s' (%s) exists in input", n.config.Model.Attributes[attributeId].Key, attributeId)
			existenceResults = append(existenceResults, true)
		} else {
			n.logger.Info().Ctx(ioData.Ctx).Msgf("%s", failureReason)
			existenceResults = append(existenceResults, false)
		}
	}

	// Memoize the result
	pass := true
	for _, result := range existenceResults {
		if !result {
			pass = false
			break
		}
	}

	// Return a new IoData instance containing the output
	return &node.IoData{
		Data: dataBytes,
		Type: node.BytesDataType, // Set the type as BytesDataType
		Ctx:  ioData.Ctx,
	}, pass, false, nil
}
