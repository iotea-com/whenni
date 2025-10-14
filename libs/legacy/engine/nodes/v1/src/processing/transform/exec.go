package main

import (
	"bytes"
	"fmt"

	"github.com/goccy/go-json"

	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
)

func (n *TransformProcessingNode) Exec(params node.ExecParams) node.Error {
Loop:
	for {
		select {
		case <-params.Ctx.Done():
			break Loop
		case ioData, ok := <-n.inputChannels[0].Channel:
			if !ok {
				break Loop
			}

			output, invalid, err := n.handle(&ioData)
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
				// No error and input was valid, so we can continue the channel workflow
				n.outputChannels[0].Channel <- *output

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

func (n *TransformProcessingNode) handle(ioData *node.IoData) (*node.IoData, bool, error) {
	defer func() {
		if r := recover(); r != nil {
			n.logger.Error().Msgf("Recovered from panic: %v", r)
		}
	}()

	// Ensure the input type is correct
	dataBytes, ok := ioData.Data.([]byte)
	if !ok {
		return nil, true, fmt.Errorf("input data has wrong type, expected bytes")
	}

	n.logger.Info().Ctx(ioData.Ctx).Str("data", string(dataBytes)).Msg("Input data")

	// Initialize the input map
	input := make(map[string]json.RawMessage)

	// Use a JSON decoder
	decoder := json.NewDecoder(bytes.NewReader(dataBytes))
	decoder.DisallowUnknownFields() // Disallow unknown fields

	// Unpack the input into a JSON struct
	err := decoder.Decode(&input)
	if err != nil {
		return nil, true, fmt.Errorf("failed to unmarshal byte data: %v, input data: %v", err, string(dataBytes))
	}

	// Transform the input data to the output data
	outputData, maskedOutputData, err := n.config.ModelInput.TransformData(&n.config.ModelOutput, input, n.config.Mapping)
	if err != nil {
		return nil, true, fmt.Errorf("failed to transform input data: %v", err)
	}

	// Marshal and return the output data
	outputBytes, err := json.Marshal(outputData)
	if err != nil {
		return nil, true, fmt.Errorf("failed to marshal output data: %v", err)
	}

	// Marshal and return the masked output data
	maskedOutputBytes, err := json.Marshal(maskedOutputData)
	if err != nil {
		return nil, true, fmt.Errorf("failed to marshal masked output data: %v", err)
	}

	n.logger.Info().Ctx(ioData.Ctx).Str("data", string(maskedOutputBytes)).Msg("Output data")

	return &node.IoData{
		Data: outputBytes,
		Type: node.BytesDataType,
		Ctx:  ioData.Ctx,
	}, false, nil
}
