package template

import (
	"encoding/json"
	"fmt"

	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
)

// Convert node.IoData to a map that the template engine can use
func ConvertIoDataToMap(ioData node.IoData) (map[string]any, error) {
	dataMap := make(map[string]any)
	switch v := ioData.Data.(type) {
	case map[string]any:
		dataMap = v
	case []byte:
		// Unmarshal if IoData is in JSON format
		if err := json.Unmarshal(v, &dataMap); err != nil {
			return nil, fmt.Errorf("failed to unmarshal IoData to map: %v", err)
		}
	case nil:
		return nil, nil
	default:
		return nil, fmt.Errorf("invalid data type in IoData: %T", ioData.Data)
	}

	return dataMap, nil
}
