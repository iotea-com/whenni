package channelsValidate

import (
	"fmt"

	gruenthttp "github.com/ongruent/gruent/libs/http"
	"github.com/ongruent/gruent/libs/legacy/engine/channels"
	resolveModels "github.com/ongruent/gruent/libs/legacy/engine/dependencies/models/resolve"
	resolveThings "github.com/ongruent/gruent/libs/legacy/engine/dependencies/things/resolve"
	"github.com/ongruent/gruent/services/http-api/services/sqlc"
	"go.opentelemetry.io/otel/attribute"
)

func execute(request *gruenthttp.Request[Input]) (*Output, error) {
	request.Span.AddEvent("execute")
	request.Span.SetAttributes(
		attribute.String("request.Input.SpaceId", request.Input.SpaceId),
	)

	// Create validation errors
	validationErrors := &channels.ValidationError{
		ChannelErrors: []string{},
		NodeErrors:    map[string][]string{},
	}

	// Expand the channel config
	for i, node := range request.Input.Config.Nodes {
		modelsResolveConfig := resolveModels.ResolveConfig{
			Metadata:    &node.Metadata,
			SqlcQueries: sqlc.Queries,
		}

		if err := resolveModels.ResolveInNode(&modelsResolveConfig); err != nil {
			validationErrors.NodeErrors[node.Id] = append(validationErrors.NodeErrors[node.Id], fmt.Sprintf("could not expand models in node %v, %v", node.Metadata.Name, err))
		}

		thingsExpandConfig := resolveThings.ExpandConfig{
			Metadata:    &node.Metadata,
			SqlcQueries: sqlc.Queries,
		}

		if err := resolveThings.ExpandInNode(&thingsExpandConfig); err != nil {
			validationErrors.NodeErrors[node.Id] = append(validationErrors.NodeErrors[node.Id], fmt.Sprintf("could not expand things in node %v, %v", node.Metadata.Name, err))
		}

		request.Input.Config.Nodes[i] = node
	}

	// Validate the expanded channel config
	ve := request.Input.Config.Validate()

	// Combined validation errors
	if ve != nil {
		validationErrors.ChannelErrors = append(validationErrors.ChannelErrors, ve.ChannelErrors...)
		validationErrors.NodeErrors = mergeNodeErrors(validationErrors.NodeErrors, ve.NodeErrors)
	}

	output := &Output{
		ChannelErrors: validationErrors.ChannelErrors,
		NodeErrors:    validationErrors.NodeErrors,
	}

	return output, nil
}

func mergeNodeErrors(existingErrors map[string][]string, newErrors map[string][]string) map[string][]string {
	for nodeId, errors := range newErrors {
		if existingErrors[nodeId] == nil {
			existingErrors[nodeId] = []string{}
		}
		existingErrors[nodeId] = append(existingErrors[nodeId], errors...)
	}
	return existingErrors
}
