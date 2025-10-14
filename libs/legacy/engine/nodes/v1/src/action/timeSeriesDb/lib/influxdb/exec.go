package influxdbNode

import (
	"context"
	"fmt"
	"time"

	"github.com/goccy/go-json"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
)

func (n *InfluxDBActionNode) Exec(params node.ExecParams) node.Error {
Loop:
	for {
		select {
		case <-params.Ctx.Done():
			break Loop
		case ioData := <-n.inputChannels[0].Channel:
			invalid, err := n.processData(ioData)
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

func (n *InfluxDBActionNode) processData(ioData node.IoData) (bool, error) {
	var dataRaw json.RawMessage
	var err error

	switch ioData.Data.(type) {
	case []byte:
		// If it's already a byte slice, assume it's JSON
		dataRaw = ioData.Data.([]byte)

	case string:
		// If it's a string, convert to byte slice
		dataRaw = []byte(ioData.Data.(string))

	case map[string]any:
		// If it's a map, marshal to JSON
		dataRaw, err = json.Marshal(ioData.Data)
		if err != nil {
			return true, fmt.Errorf("failed to marshal input data: %v", err)
		}

	default:
		return true, fmt.Errorf("unexpected data type in IoData: %T", ioData.Data)
	}

	// Attempt to write
	if err := n.writeData(ioData.Ctx, dataRaw); err != nil {
		return false, fmt.Errorf("unable to write to database: %v", err)
	}

	return false, nil
}

func (n *InfluxDBActionNode) writeData(dataCtx context.Context, data json.RawMessage) error {
	// Only tag for now is a test user
	tags := map[string]string{
		"user": "test_user",
	}

	// Convert data to map[string]json.RawMessage
	var dataMap map[string]json.RawMessage
	if err := json.Unmarshal(data, &dataMap); err != nil {
		return fmt.Errorf("failed to unmarshal input data: %v", err)
	}

	// Convert json.RawMessage to actual data values and populate fields
	fields := make(map[string]any)
	for key, rawValue := range dataMap {
		var value any
		if err := json.Unmarshal(rawValue, &value); err != nil {
			return fmt.Errorf("error unmarshalling raw JSON for key %s: %v", key, err)
		}
		fields[key] = value
	}

	// Create a new point with the measurement name, tags, fields, and the current time
	point := influxdb2.NewPoint("measurements", tags, fields, time.Now())

	// Prepare the log data
	logData := map[string]any{
		"measurement": "measurements",
		"tags":        tags,
		"fields":      fields,
		"time":        point.Time(),
	}

	logDataBytes, err := json.Marshal(logData)
	if err != nil {
		return fmt.Errorf("error marshalling log data: %v", err)
	}

	n.logger.Info().Ctx(dataCtx).Str("data", string(logDataBytes)).Msgf("Writing point to bucket %v", n.config.Bucket)

	// Write the point using a blocking API to ensure completion
	if err := n.writeAPI.WritePoint(context.Background(), point); err != nil {
		return fmt.Errorf("failed to write point to InfluxDB: %v", err)
	}

	// Log the relevant values of the point
	n.logger.Info().Ctx(dataCtx).Msgf("Point written to bucket %v", n.config.Bucket)

	return nil
}
