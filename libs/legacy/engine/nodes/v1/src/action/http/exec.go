package main

import (
	"bytes"
	"io"
	"net/http"

	"github.com/goccy/go-json"

	"fmt"

	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
)

func (n *HttpActionNode) Exec(params node.ExecParams) node.Error {
Loop:
	for {
		select {
		case <-params.Ctx.Done():
			break Loop
		case ioData := <-n.inputChannels[0].Channel:
			invalid, err := n.handlePayload(ioData)
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

func (n *HttpActionNode) handlePayload(ioData node.IoData) (bool, error) {
	var payloadDataReader io.Reader
	// If there is a request body
	if ioData.Data != nil {
		payloadData, ok := ioData.Data.([]byte)
		if !ok {
			return true, fmt.Errorf("expected data of type []byte, got %T", ioData.Data)
		}

		// Check if payloadData is already JSON; if so, no need to marshal again
		if !json.Valid(payloadData) {
			pd, err := json.Marshal(payloadData)
			if err != nil {
				return false, fmt.Errorf("error creating request: %v", err)
			}
			payloadData = pd
		}

		payloadDataReader = bytes.NewReader(payloadData)
	}

	// Prevent HTTP bodies in GET or HEAD requests
	if n.config.Method == "GET" || n.config.Method == "HEAD" {
		payloadDataReader = nil
	}

	// Create the request URL
	url := fmt.Sprintf(
		"%s://%s:%d%s",
		n.config.ThingHttpServer.Protocol,
		n.config.ThingHttpServer.Host,
		n.config.ThingHttpServer.Port,
		n.config.Path,
	)

	// Create the request
	request, err := http.NewRequest(n.config.Method, url, payloadDataReader)
	if err != nil {
		return false, fmt.Errorf("error creating response: %v", err)
	}

	if payloadDataReader != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	n.logger.Info().Ctx(ioData.Ctx).Msgf("Sending %v HTTP request to %s", n.config.Method, url)

	// Send the request
	response, err := n.client.Do(request)
	if err != nil {
		return false, fmt.Errorf("error receiving response: %v", err)
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return false, fmt.Errorf("error parsing response: %v", err)
	}

	result := map[string]any{
		"request": map[string]any{
			"method":  n.config.Method,
			"url":     url,
			"body":    n.config.Body,
			"headers": n.config.Headers,
		},
		"response": map[string]any{
			"status": response.StatusCode,
			"body":   string(responseBody),
		},
	}

	resultBytes, err := json.Marshal(result)
	if err != nil {
		return false, fmt.Errorf("error marshalling response: %v", err)
	}

	n.logger.Info().Ctx(ioData.Ctx).Str("data", string(resultBytes)).Msg("Reponse received")

	outputData := node.IoData{
		Data: resultBytes,
		Type: node.MapDataType,
		Ctx:  ioData.Ctx,
	}

	select {
	case n.outputChannels[0].Channel <- outputData:
	default:
		return false, nil
	}

	return false, nil
}
