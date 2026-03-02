package channels

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/go-playground/validator/v10"

	httpActionNodeConfig "github.com/iotea-com/iotea/libs/legacy/engine/nodes/v1/src/action/http/config"
	mqttActionNodeConfig "github.com/iotea-com/iotea/libs/legacy/engine/nodes/v1/src/action/mqtt/config"
	booleanConditionalNodeConfig "github.com/iotea-com/iotea/libs/legacy/engine/nodes/v1/src/conditional/boolean/config"
	existenceConditionalNodeConfig "github.com/iotea-com/iotea/libs/legacy/engine/nodes/v1/src/conditional/existence/config"
	stringCompareConditionalNodeConfig "github.com/iotea-com/iotea/libs/legacy/engine/nodes/v1/src/conditional/stringCompare/config"
	thresholdConditionalNodeConfig "github.com/iotea-com/iotea/libs/legacy/engine/nodes/v1/src/conditional/threshold/config"
	transformProcessingNodeConfig "github.com/iotea-com/iotea/libs/legacy/engine/nodes/v1/src/processing/transform/config"
	mqttSourceNodeConfig "github.com/iotea-com/iotea/libs/legacy/engine/nodes/v1/src/source/mqtt/config"
	timerSourceNodeConfig "github.com/iotea-com/iotea/libs/legacy/engine/nodes/v1/src/source/timer/config"

	"github.com/iotea-com/iotea/libs/val"
)

type ValidationError struct {
	ChannelErrors []string            `json:"channelErrors"`
	NodeErrors    map[string][]string `json:"nodeErrors"`
}

/*
Validates the channel DAG according to the channel rules.

- Only one source node may be present.

- Must have more than one node.

- Any node may only have one connected input.

- Any node may have many connected outputs, which we refer to as "branching" or "branches".

- A branch may end in any type of node (except a source node).

- No "hanging nodes", or nodes without at least one edge.

- No duplicate edges.

- Edges cannot connect inputs and outputs of the same node.

- Runtime size must be valid.
*/
func (c *Channel) Validate() *ValidationError {
	validationErrors := &ValidationError{
		ChannelErrors: []string{},
		NodeErrors:    map[string][]string{},
	}

	// Check if channel has more than one node
	if len(c.Nodes) < 2 {
		validationErrors.ChannelErrors = append(validationErrors.ChannelErrors, "Channel must have more than one node.")
	}

	// Populate nodeMap and initialize inDegree and outDegree maps
	nodeMap := make(map[string]Node)
	edgeMap := make(map[string]bool)
	inDegree := make(map[string]int)
	outDegree := make(map[string]int)
	sourceNodeCount := 0

	for _, node := range c.Nodes {
		nodeMap[node.Id] = node
		inDegree[node.Id] = 0
		outDegree[node.Id] = 0
		if node.Metadata.Type == "source" {
			sourceNodeCount++
		}
	}

	// Ensure only one source node is present
	if sourceNodeCount > 1 {
		validationErrors.ChannelErrors = append(validationErrors.ChannelErrors, "Only one source node is allowed.")
	}

	// Calculate inDegree and outDegree for each node based on edges
	for _, edge := range c.Edges {
		if _, exists := nodeMap[edge.From.NodeId]; !exists {
			validationErrors.ChannelErrors = append(validationErrors.ChannelErrors, fmt.Sprintf("Edge from node %s does not exist.", edge.From.NodeId))
		}
		if _, exists := nodeMap[edge.To.NodeId]; !exists {
			validationErrors.ChannelErrors = append(validationErrors.ChannelErrors, fmt.Sprintf("Edge to node %s does not exist.", edge.To.NodeId))
		}

		// Check for duplicate edges
		edgeKey := fmt.Sprintf("%s-%s-%s-%s", edge.From.NodeId, edge.From.IoId, edge.To.NodeId, edge.To.IoId)
		if edgeMap[edgeKey] {
			validationErrors.ChannelErrors = append(validationErrors.ChannelErrors, fmt.Sprintf("Duplicate edge detected between nodes %s and %s with IO points %s and %s.", edge.From.NodeId, edge.To.NodeId, edge.From.IoId, edge.To.IoId))
		}
		edgeMap[edgeKey] = true

		inDegree[edge.To.NodeId]++
		outDegree[edge.From.NodeId]++
	}

	// Validate no "hanging" nodes (nodes with no edges)
	for nodeId, node := range nodeMap {
		if inDegree[nodeId] == 0 && outDegree[nodeId] == 0 {
			validationErrors.ChannelErrors = append(validationErrors.ChannelErrors, fmt.Sprintf("Node %s (%s) is a hanging node with no edges.", nodeId, node.Metadata.Name))
		}
	}

	// Validate that each node has exactly one input (inDegree == 1) and at least one output
	for nodeId, degree := range inDegree {
		node := nodeMap[nodeId]

		if degree > 1 {
			validationErrors.ChannelErrors = append(validationErrors.ChannelErrors, fmt.Sprintf("Node %s (%s) has more than one input.", nodeId, node.Metadata.Name))
		}
		if degree == 0 && nodeMap[nodeId].Metadata.Type != "source" {
			validationErrors.ChannelErrors = append(validationErrors.ChannelErrors, fmt.Sprintf("Node %s (%s) has no input and is not a source node.", nodeId, node.Metadata.Name))
		}
	}

	// Validate that each edge is not connecting two inputs or two outputs
	for _, edge := range c.Edges {
		fromIoDirection := "output"
		toIoDirection := "output"

		if strings.Contains(strings.ToLower(edge.From.IoId), "input") {
			fromIoDirection = "input"
		}

		if strings.Contains(strings.ToLower(edge.To.IoId), "input") {
			toIoDirection = "input"
		}

		if fromIoDirection == "output" && toIoDirection == "output" {
			fromNode := nodeMap[edge.From.NodeId]
			toNode := nodeMap[edge.To.NodeId]
			validationErrors.ChannelErrors = append(validationErrors.ChannelErrors, fmt.Sprintf("Edge %s (between %s and %s) connects two outputs.", edge.Id, fromNode.Metadata.Name, toNode.Metadata.Name))
		}

		if fromIoDirection == "input" && toIoDirection == "input" {
			fromNode := nodeMap[edge.From.NodeId]
			toNode := nodeMap[edge.To.NodeId]
			validationErrors.ChannelErrors = append(validationErrors.ChannelErrors, fmt.Sprintf("Edge %s (between %s and %s) connects two inputs.", edge.Id, fromNode.Metadata.Name, toNode.Metadata.Name))
		}

		if edge.From.NodeId == edge.To.NodeId {
			fromNode := nodeMap[edge.From.NodeId]
			toNode := nodeMap[edge.To.NodeId]
			validationErrors.ChannelErrors = append(validationErrors.ChannelErrors, fmt.Sprintf("Edge %s (between %s and %s) connects an input to an output from the same node %s.", edge.Id, fromNode.Metadata.Name, toNode.Metadata.Name, edge.From.NodeId))
		}
	}

	// Validate runtime size
	validSizes := []string{"small", "medium", "large"}
	if !slices.Contains(validSizes, c.Runtime.Size) {
		validationErrors.ChannelErrors = append(validationErrors.ChannelErrors, fmt.Sprintf("Invalid runtime size. Must be one of: %s. Received: %s.", strings.Join(validSizes, ", "), c.Runtime.Size))
	}

	// Validate node config
	for _, node := range c.Nodes {
		nodeValidationErrors := node.ValidateConfig()
		if len(nodeValidationErrors) > 0 {
			validationErrors.NodeErrors[node.Id] = nodeValidationErrors
		}
	}

	// If there are no validation errors, return nil
	if len(validationErrors.ChannelErrors) <= 0 && len(validationErrors.NodeErrors) <= 0 {
		return nil
	}

	return validationErrors
}

func (n *Node) ValidateConfig() []string {
	validationErrors := []string{}

	v := validator.New()
	v.RegisterValidation("valid_paths", val.IsValidPaths)
	v.RegisterValidation("mqtt_topic", val.MqttTopic)

	// Detect the type of node
	nodeType := n.Metadata.Type
	switch nodeType {
	case "source":
		switch n.Metadata.Label {
		case "http":
			// The attributes for the HTTP source node are populated
			// by the resolve package in the orchestrator service.
			// We can't validate here because the attributes don't exist yet.
			break
		case "mqtt":
			// Convert the generic node config to the mqtt node config
			mqttSourceNodeConfig := mqttSourceNodeConfig.MqttSourceNodeConfig{}
			err := json.Unmarshal(n.Metadata.Config, &mqttSourceNodeConfig)
			if err != nil {
				validationErrors = append(validationErrors, fmt.Sprintf("Invalid config for node %s: %s.", n.Metadata.Name, err.Error()))
			}

			// Validate the mqtt node config
			err = v.Struct(mqttSourceNodeConfig)
			if err != nil {
				for _, err := range err.(validator.ValidationErrors) {
					validationErrors = append(validationErrors, fmt.Sprintf("Invalid config for node %s: %s.", n.Metadata.Name, err.Error()))
				}
			}
		case "messageQueue":
			// TODO: implement subnode validation
			break
		case "timer":
			// Convert the generic node config to the timer node config
			timerSourceNodeConfig := timerSourceNodeConfig.TimerSourceNodeConfig{}
			err := json.Unmarshal(n.Metadata.Config, &timerSourceNodeConfig)
			if err != nil {
				validationErrors = append(validationErrors, fmt.Sprintf("Invalid config for node %s: %s.", n.Metadata.Name, err.Error()))
			}

			// Validate the timer node config
			err = v.Struct(timerSourceNodeConfig)
			if err != nil {
				for _, err := range err.(validator.ValidationErrors) {
					validationErrors = append(validationErrors, fmt.Sprintf("Invalid config for node %s: %s.", n.Metadata.Name, err.Error()))
				}
			}
		}
	case "processing":
		switch n.Metadata.Label {
		case "transform":
			// Convert the generic node config to the transform node config
			transformProcessingNodeConfig := transformProcessingNodeConfig.TransformProcessingNodeConfig{}
			err := json.Unmarshal(n.Metadata.Config, &transformProcessingNodeConfig)
			if err != nil {
				validationErrors = append(validationErrors, fmt.Sprintf("Invalid config for node %s: %s.", n.Metadata.Name, err.Error()))
			}

			// Validate the transform node config
			err = v.Struct(transformProcessingNodeConfig)
			if err != nil {
				for _, err := range err.(validator.ValidationErrors) {
					validationErrors = append(validationErrors, fmt.Sprintf("Invalid config for node %s: %s.", n.Metadata.Name, err.Error()))
				}
			}
		}
	case "conditional":
		switch n.Metadata.Label {
		case "threshold":
			// Convert the generic node config to the threshold node config
			thresholdConditionalNodeConfig := thresholdConditionalNodeConfig.ThresholdConditionalNodeConfig{}
			err := json.Unmarshal(n.Metadata.Config, &thresholdConditionalNodeConfig)
			if err != nil {
				validationErrors = append(validationErrors, fmt.Sprintf("Invalid config for node %s: %s.", n.Metadata.Name, err.Error()))
			}

			// Validate the threshold node config
			err = v.Struct(thresholdConditionalNodeConfig)
			if err != nil {
				for _, err := range err.(validator.ValidationErrors) {
					validationErrors = append(validationErrors, fmt.Sprintf("Invalid config for node %s: %s.", n.Metadata.Name, err.Error()))
				}
			}
		case "existence":
			// Convert the generic node config to the existence node config
			existenceConditionalNodeConfig := existenceConditionalNodeConfig.ExistenceConditionalNodeConfig{}
			err := json.Unmarshal(n.Metadata.Config, &existenceConditionalNodeConfig)
			if err != nil {
				validationErrors = append(validationErrors, fmt.Sprintf("Invalid config for node %s: %s.", n.Metadata.Name, err.Error()))
			}

			// Validate the existence node config
			err = v.Struct(existenceConditionalNodeConfig)
			if err != nil {
				for _, err := range err.(validator.ValidationErrors) {
					validationErrors = append(validationErrors, fmt.Sprintf("Invalid config for node %s: %s.", n.Metadata.Name, err.Error()))
				}
			}
		case "stringCompare":
			// Convert the generic node config to the string compare node config
			stringCompareConditionalNodeConfig := stringCompareConditionalNodeConfig.StringCompareConditionalNodeConfig{}
			err := json.Unmarshal(n.Metadata.Config, &stringCompareConditionalNodeConfig)
			if err != nil {
				validationErrors = append(validationErrors, fmt.Sprintf("Invalid config for node %s: %s.", n.Metadata.Name, err.Error()))
			}

			// Validate the string compare node config
			err = v.Struct(stringCompareConditionalNodeConfig)
			if err != nil {
				for _, err := range err.(validator.ValidationErrors) {
					validationErrors = append(validationErrors, fmt.Sprintf("Invalid config for node %s: %s.", n.Metadata.Name, err.Error()))
				}
			}
		case "boolean":
			// Convert the generic node config to the boolean node config
			booleanConditionalNodeConfig := booleanConditionalNodeConfig.BooleanConditionalNodeConfig{}
			err := json.Unmarshal(n.Metadata.Config, &booleanConditionalNodeConfig)
			if err != nil {
				validationErrors = append(validationErrors, fmt.Sprintf("Invalid config for node %s: %s.", n.Metadata.Name, err.Error()))
			}

			// Validate the boolean node config
			err = v.Struct(booleanConditionalNodeConfig)
			if err != nil {
				validationErrors = append(validationErrors, fmt.Sprintf("Invalid config for node %s: %s.", n.Metadata.Name, err.Error()))
			}

			// Validate the boolean node config
			err = v.Struct(booleanConditionalNodeConfig)
			if err != nil {
				for _, err := range err.(validator.ValidationErrors) {
					validationErrors = append(validationErrors, fmt.Sprintf("Invalid config for node %s: %s.", n.Metadata.Name, err.Error()))
				}
			}
		}
	case "action":
		switch n.Metadata.Label {
		case "http":
			// Convert the generic node config to the http action node config
			httpActionNodeConfig := httpActionNodeConfig.HttpActionNodeConfig{}
			err := json.Unmarshal(n.Metadata.Config, &httpActionNodeConfig)
			if err != nil {
				validationErrors = append(validationErrors, fmt.Sprintf("Invalid config for node %s: %s.", n.Metadata.Name, err.Error()))
			}

			// Validate the http node config
			err = v.Struct(httpActionNodeConfig)
			if err != nil {
				for _, err := range err.(validator.ValidationErrors) {
					validationErrors = append(validationErrors, fmt.Sprintf("Invalid config for node %s: %s.", n.Metadata.Name, err.Error()))
				}
			}
		case "mqtt":
			// Convert the generic node config to the mqtt action node config
			mqttActionNodeConfig := mqttActionNodeConfig.MqttActionNodeConfig{}
			err := json.Unmarshal(n.Metadata.Config, &mqttActionNodeConfig)
			if err != nil {
				validationErrors = append(validationErrors, fmt.Sprintf("Invalid config for node %s: %s.", n.Metadata.Name, err.Error()))
			}

			// Validate the mqtt node config
			err = v.Struct(mqttActionNodeConfig)
			if err != nil {
				for _, err := range err.(validator.ValidationErrors) {
					validationErrors = append(validationErrors, fmt.Sprintf("Invalid config for node %s: %s.", n.Metadata.Name, err.Error()))
				}
			}
		case "messageQueue":
			// TODO: implement subnode validation
			break
		case "fileStorage":
			// TODO: implement subnode validation
			break
		case "timeSeriesDb":
			// TODO: implement subnode validation
			break
		case "notification":
			// TODO: implement subnode validation
			break
		case "documentDb":
			// TODO: implement subnode validation
			break
		}
	default:
		validationErrors = append(validationErrors, fmt.Sprintf("Invalid node type: %s.", nodeType))
	}

	return validationErrors
}
