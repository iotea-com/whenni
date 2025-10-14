package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/iotea-com/iotea/libs/engine/logs"
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
)

func (n *MetricActionNode) Exec(params node.ExecParams) node.Error {
Loop:
	for {
		select {
		case <-params.Ctx.Done():
			break Loop
		case ioData := <-n.inputChannels[0].Channel:
			invalid, err := n.handlePayload(ioData)
			if err != nil {
				if invalid {
					// Safely notify the runtime that data was invalid
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

func (n *MetricActionNode) handlePayload(ioData node.IoData) (bool, error) {
	if ioData.Data == nil {
		return false, fmt.Errorf("input data is nil")
	}

	// Parse the input data
	var inputData map[string]json.RawMessage
	if err := json.Unmarshal(ioData.Data.([]byte), &inputData); err != nil {
		return false, fmt.Errorf("failed to parse input data as JSON: %v", err)
	}

	// Validate against model
	if err := n.config.Model.Compare(inputData); err != nil {
		return true, fmt.Errorf("input data does not match model: %v", err)
	}

	// Extract value and metadata based on configured fields
	var metricValue float64
	var metadataValue interface{}

	// Get the value field
	valueFieldId := n.config.Attribute
	valueAttr, exists := n.config.Model.Attributes[valueFieldId]
	if !exists {
		return false, fmt.Errorf("value field ID '%s' not found in model attributes", valueFieldId)
	}

	// Look for the field in the input data using the attribute's key
	valueData, exists := inputData[valueAttr.Key]
	if !exists {
		return true, fmt.Errorf("missing required value field '%s' in payload", valueAttr.Key)
	}
	if err := json.Unmarshal(valueData, &metricValue); err != nil {
		return true, fmt.Errorf("failed to unmarshal number value for '%s': %v", valueAttr.Key, err)
	}

	// Get the metadata field
	metadataFieldId := n.config.Metadata
	metadataAttr, exists := n.config.Model.Attributes[metadataFieldId]
	if !exists {
		return false, fmt.Errorf("metadata field ID '%s' not found in model attributes", metadataFieldId)
	}

	// Look for the field in the input data using the attribute's key
	metadataData, exists := inputData[metadataAttr.Key]
	if !exists {
		return true, fmt.Errorf("missing required metadata field '%s' in payload", metadataAttr.Key)
	}

	// Unmarshal metadata based on its type
	switch metadataAttr.Type {
	case "string":
		var strValue string
		if err := json.Unmarshal(metadataData, &strValue); err != nil {
			return true, fmt.Errorf("failed to unmarshal string metadata for '%s': %v", metadataAttr.Key, err)
		}
		metadataValue = strValue
	case "number":
		var numValue float64
		if err := json.Unmarshal(metadataData, &numValue); err != nil {
			return true, fmt.Errorf("failed to unmarshal number metadata for '%s': %v", metadataAttr.Key, err)
		}
		metadataValue = numValue
	case "boolean":
		var boolValue bool
		if err := json.Unmarshal(metadataData, &boolValue); err != nil {
			return true, fmt.Errorf("failed to unmarshal boolean metadata for '%s': %v", metadataAttr.Key, err)
		}
		metadataValue = boolValue
	default:
		return false, fmt.Errorf("unsupported metadata type '%s' for field '%s'", metadataAttr.Type, metadataAttr.Key)
	}

	// Create metadata map
	metadata := map[string]interface{}{
		metadataAttr.Key: metadataValue,
	}

	// Insert into ClickHouse
	ctx := context.Background()
	query := `
		INSERT INTO custom_metrics (
			Key,
			Metadata,
			SpaceId,
			ChannelId,
			Value,
			Timestamp
		) VALUES (
			?,
			?,
			?,
			?,
			?,
			?
		)
	`

	// Get space and channel IDs from context
	spaceId, ok := ioData.Ctx.Value(logs.SpaceIDKey).(string)
	if !ok {
		return false, fmt.Errorf("spaceId not found in context")
	}
	channelId, ok := ioData.Ctx.Value(logs.ChannelIDKey).(string)
	if !ok {
		return false, fmt.Errorf("channelId not found in context")
	}

	err := (*n.clickhouseClient).Exec(ctx, query,
		valueAttr.Key,
		metadata,
		spaceId,
		channelId,
		metricValue,
		time.Now().UTC(),
	)

	if err != nil {
		return false, fmt.Errorf("failed to insert metric into ClickHouse: %v", err)
	}

	n.logger.Info().Ctx(ioData.Ctx).Msgf("Saved metric: metadata=%v, key=%s, value=%f", metadata, valueAttr.Key, metricValue)

	return false, nil
}
