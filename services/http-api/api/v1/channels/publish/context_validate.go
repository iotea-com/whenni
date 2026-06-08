package channelsPublish

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	ioteapermissions "github.com/iotea-com/iotea/libs/http/permissions"
	"github.com/iotea-com/iotea/libs/legacy/engine/channels"
	resolveModels "github.com/iotea-com/iotea/libs/legacy/engine/dependencies/models/resolve"
	resolveThings "github.com/iotea-com/iotea/libs/legacy/engine/dependencies/things/resolve"
	"github.com/iotea-com/iotea/services/http-api/config"
	"github.com/iotea-com/iotea/services/http-api/services/sqlc"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func contextValidate(request *ioteahttp.Request[Input]) error {
	request.Span.AddEvent("contextValidate")

	authorizeRequestParams := ioteahttp.AuthorizeRequestParams{
		BearerToken: request.Input.BearerToken,
		JwtSecret:   config.VaultConf.JwtSecret,
		ScopeId:     request.Input.SpaceId,
		Namespace:   ioteapermissions.NamespaceChannels,
		Action:      ioteapermissions.ActionUpdate,
	}

	authorizeErr := request.Authorize(authorizeRequestParams)
	if authorizeErr != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "authorize"),
			attribute.String("error.message", authorizeErr.Error()),
			attribute.String("context_validation.status", "fail"),
		)
		return authorizeErr
	}

	// Get the channel config from the database
	dbCtx, dbSpan := otel.Tracer("sqlc").Start(request.Context, "Get channel")
	channelRecord, err := sqlc.Queries.GetChannel(dbCtx, request.Input.ChannelId)
	if err != nil {
		if err == pgx.ErrNoRows {
			message := fmt.Sprintf("no channel found with ID %s", request.Input.ChannelId)
			dbSpan.SetAttributes(
				attribute.String("error.type", "database"),
				attribute.String("error.message", message),
				attribute.String("context_validation.status", "fail"),
			)
			dbSpan.End()

			errorResponse := ioteahttp.NewErrorResponse([]any{message})
			responseJson, _ := errorResponse.MarshalJson()
			return fiber.NewError(fiber.StatusBadRequest, string(responseJson))
		}

		message := fmt.Sprintf("error retrieving channel from the database: %s", err)
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", message),
			attribute.String("context_validation.status", "fail"),
		)
		dbSpan.End()

		errorResponse := ioteahttp.NewErrorResponse([]any{message})
		responseJson, _ := errorResponse.MarshalJson()
		return fiber.NewError(fiber.StatusBadRequest, string(responseJson))
	}

	dbSpan.End()

	var c = &channels.Channel{}
	err = json.Unmarshal(channelRecord.Config, c)
	if err != nil {
		errorResponse := ioteahttp.NewErrorResponse([]any{err.Error()})
		responseJson, _ := errorResponse.MarshalJson()
		return fiber.NewError(fiber.StatusBadRequest, string(responseJson))
	}

	// Create validation errors
	validationErrors := &channels.ValidationError{
		ChannelErrors: []string{},
		NodeErrors:    map[string][]string{},
	}

	// Expand the channel config
	for i, node := range c.Nodes {
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

		c.Nodes[i] = node
	}

	// Validate the expanded channel config
	ve := c.Validate()

	// Combined validation errors
	if ve != nil {
		validationErrors.ChannelErrors = append(validationErrors.ChannelErrors, ve.ChannelErrors...)
		validationErrors.NodeErrors = mergeNodeErrors(validationErrors.NodeErrors, ve.NodeErrors)
	}

	// Combine into a single slice
	var combinedErrors []string
	for _, nodeErrors := range validationErrors.NodeErrors {
		combinedErrors = append(combinedErrors, nodeErrors...)
	}
	combinedErrors = append(combinedErrors, validationErrors.ChannelErrors...)

	if len(combinedErrors) > 0 {
		errorsAny := make([]any, len(combinedErrors))
		for i, err := range combinedErrors {
			errorsAny[i] = err
		}
		errorResponse := ioteahttp.NewErrorResponse(errorsAny)
		responseJson, err := errorResponse.MarshalJson()
		if err != nil {
			request.Span.SetAttributes(
				attribute.String("error.type", "validate_channel"),
				attribute.String("error.message", err.Error()),
				attribute.String("context_validation.status", "fail"),
			)
			errorResponse = ioteahttp.NewErrorResponse([]any{err.Error()})
			responseJson, _ = errorResponse.MarshalJson()
			return fiber.NewError(fiber.StatusBadRequest, string(responseJson))
		}
		return fiber.NewError(fiber.StatusBadRequest, string(responseJson))
	}

	request.Span.SetAttributes(attribute.String("context_validation.status", "pass"))
	return nil
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
