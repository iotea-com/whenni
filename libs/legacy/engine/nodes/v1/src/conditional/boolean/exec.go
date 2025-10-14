package main

import (
	"fmt"

	"github.com/goccy/go-json"

	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
)

func (n *BooleanConditionalNode) Exec(params node.ExecParams) node.Error {
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
func (n *BooleanConditionalNode) handle(ioData node.IoData) (*node.IoData, bool, bool, error) {
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
		return nil, false, true, fmt.Errorf("failed to unmarshal byte data: %v", err)
	}

	conditionResults := []bool{}
	// Evaluate conditions and return upon the first successful one
	for i, condition := range n.config.Conditions {
		if failureReason := n.config.Model.EvaluateBooleanCondition(input, condition); failureReason == nil {
			n.logger.Info().Ctx(ioData.Ctx).Msgf("Condition %d passed", i)
			conditionResults = append(conditionResults, true)
		} else {
			n.logger.Info().Ctx(ioData.Ctx).Msgf("Condition failed: %s", failureReason)
			conditionResults = append(conditionResults, false)
		}
	}
	pass := false
	switch n.config.LogicalOperator {
	case "AND":
		pass = true
		for _, result := range conditionResults {
			if !result {
				pass = false
			}
		}

	case "OR":
		for _, result := range conditionResults {
			if result {
				pass = true
			}
		}
	}

	// Return a new IoData instance containing the output
	return &node.IoData{
		Data: dataBytes,
		Type: node.BytesDataType, // Set the type as BytesDataType
		Ctx:  ioData.Ctx,
	}, pass, false, nil
}
