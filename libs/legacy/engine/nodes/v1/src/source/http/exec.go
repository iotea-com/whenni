package main

import (
	"fmt"
	"time"

	"github.com/goccy/go-json"

	"github.com/gofiber/fiber/v2"
	"github.com/iotea-com/iotea/libs/engine/logs"
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	httpSourceNodeConfig "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/source/http/config"
)

func (n *HttpSourceNode) Exec(params node.ExecParams) node.Error {
	// Start the server
	go func() {
		err := n.server.Listen(fmt.Sprintf(":%d", n.config.ThingHttpServer.Port))
		if err != nil {
			n.errChan <- fmt.Errorf("server returned on Listen(), %v", err)
		}
	}()

Loop:
	for {
		select {
		case <-params.Ctx.Done():
			break Loop
		case err := <-n.errChan:
			if err != nil {
				return node.Error{
					Type:   node.FatalError,
					Reason: err.Error(),
				}
			}
		}
	}

	return node.Error{
		Type: node.NoError,
	}
}

// httpRequestSubscriptionCallback will be called by the exec function when a new HTTP request
// is received on the server's route. This function does not do any checks on the data,
// and simply forwards it to its output channel as a byte array.
func (n *HttpSourceNode) httpRequestSubscriptionCallback(c *fiber.Ctx) error {
	method := c.Method()
	payload := c.Body()

	// Create an unique context for the request
	dataCtx := logs.NewDataCtx()

	// Notify the runtime that data was received
	n.notifyChannel <- node.Notification{
		Type:    node.NotifyInputDataReceived,
		DataCtx: dataCtx,
	}

	n.logger.Info().Ctx(dataCtx).Str("data", string(payload)).Msg("Request received")

	// Check if the payload is JSON or binary data based on contentType
	contentType := c.Get("Content-Type")

	if len(contentType) > 0 && string(contentType[0]) == "application/json" {
		// Validate that the payload is valid JSON
		var requestBody map[string]any
		if len(payload) > 0 && method != fiber.MethodGet && method != fiber.MethodHead {
			err := json.Unmarshal(payload, &requestBody)
			if err != nil {
				// Notify the runtime that the data is invalid
				n.notifyChannel <- node.Notification{
					Type:    node.NotifyDataInvalid,
					DataCtx: dataCtx,
				}

				return fiber.NewError(fiber.StatusBadRequest, err.Error())
			}
		}
	}

	// Create the node data type
	outputData := node.IoData{
		Data: payload,
		Type: n.outputChannels[0].Type[0],
		Ctx:  dataCtx,
	}

	// Send data to the next node or the trigger
	n.outputChannels[0].Channel <- outputData

	// Notify the runtime that the data was processed
	n.notifyChannel <- node.Notification{
		Type:    node.NotifyDataProcessed,
		DataCtx: dataCtx,
	}

	// Ok now we need to wait for the response, so first, we need to get the execution ID
	// so that we can keep track of the request
	executionId := dataCtx.Value(logs.ExecutionIDKey).(string)

	// Get timeout from configuration or use default
	timeout := time.Duration(n.config.ResponseTimeout) * time.Second

	// Register this request with the router
	responseChan := n.router.RegisterRequest(executionId, timeout)
	defer n.router.UnregisterRequest(executionId)

	// Wait for response
	responseDone := make(chan error, 1)
	go func() {
		response := <-responseChan // Wait for either response or timeout

		if response.TimedOut {
			n.logger.Warn().
				Str("flag", string(logs.TimeoutFlag)).
				Ctx(dataCtx).
				Str("data", string(payload)).
				Msgf("Request timed out waiting for response body, is your logic taking longer than the timeout of %v seconds you've set?", n.config.ResponseTimeout)

			responseDone <- fiber.NewError(
				fiber.StatusGatewayTimeout,
				"request timed out waiting for response",
			)
			return
		}

		// Try to get HttpResponseInfo from the response
		httpResponse, ok := response.NodeData.Data.(httpSourceNodeConfig.HttpResponseInfo)
		if !ok {
			n.logger.Error().
				Ctx(dataCtx).
				Msg("Received invalid response type - expected HttpResponseInfo")

			// Let the request timeout instead of sending an error response
			// This ensures consistent behavior with the timeout case
			return
		}

		// Set response code
		if httpResponse.ResponseCode > 0 {
			c.Status(httpResponse.ResponseCode)
		}

		// Set headers if provided
		if len(httpResponse.Headers) > 0 {
			for key, value := range httpResponse.Headers {
				c.Set(key, value)
			}
		}

		// Send the response body
		responseDone <- c.Send(httpResponse.ResponseBody)

		n.logger.Info().Ctx(dataCtx).Str("data", string(httpResponse.ResponseBody)).Msgf("Response set successfully, code %d", httpResponse.ResponseCode)
	}()

	return <-responseDone
}
