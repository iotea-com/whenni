package router

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/iotea-com/iotea/libs/engine/logs"
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	httpSourceNodeConfig "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/source/http/config"
	"github.com/rs/zerolog"
)

type ResponseChannel struct {
	Id           string   `json:"id" validate:"required"`
	DataType     []string `json:"dataType" validate:"required"`
	AllowedEdge  []string `json:"allowedEdge"`
	ResponseCode int      `json:"responseCode"`
	Timeout      int      `json:"timeout"`
}

type RouterResponse struct {
	ChannelId string
	NodeData  node.IoData
	TimedOut  bool
}

type responseWaiter struct {
	responseChan chan RouterResponse
	expireTime   time.Time
	expired      atomic.Bool
}

type ResponseRouter struct {
	// Map of executionId to response waiter
	waiters sync.Map

	logger zerolog.Logger
}

func NewResponseRouter(logger zerolog.Logger) *ResponseRouter {
	return &ResponseRouter{
		logger: logger,
	}
}

// RegisterRequest creates a new waiter for a specific executionId
func (r *ResponseRouter) RegisterRequest(executionId string, timeout time.Duration) chan RouterResponse {
	responseChan := make(chan RouterResponse, 1)
	waiter := &responseWaiter{
		responseChan: responseChan,
		expireTime:   time.Now().Add(timeout),
	}
	r.waiters.Store(executionId, waiter)

	// Start timeout goroutine
	go func() {
		time.Sleep(timeout)
		if waiter.expired.CompareAndSwap(false, true) {
			// Only send timeout if no response was received
			responseChan <- RouterResponse{
				TimedOut: true,
			}
		}
		// Note: We don't unregister here - that's still handled by the caller
	}()

	return responseChan
}

// UnregisterRequest removes a waiter for a specific executionId
func (r *ResponseRouter) UnregisterRequest(executionId string) {
	if waiter, ok := r.waiters.LoadAndDelete(executionId); ok {
		// The waiter is a *responseWaiter, not a channel
		w := waiter.(*responseWaiter)
		close(w.responseChan) // Close the response channel instead
	}
}

// RouteResponse routes the response data to the correct request handler
func (r *ResponseRouter) RouteResponse(channelId string, data node.IoData) error {
	// First thing we need to check is if the data passed is of type HttpResponseInfo
	_, ok := data.Data.(httpSourceNodeConfig.HttpResponseInfo)
	if !ok {
		return fmt.Errorf("Data is not of type HttpResponseInfo, meaning a node that was not of type HttpActionResponseNode sent a response")
	}

	// Check that the context has an execution ID before attempting to get it
	executionIdValue := data.Ctx.Value(logs.ExecutionIDKey)
	if executionIdValue == nil {
		return fmt.Errorf("no execution ID found in context")
	}

	executionId, ok := executionIdValue.(string)
	if !ok {
		return fmt.Errorf("invalid execution ID type in context, expected string")
	}

	if waiter, ok := r.waiters.Load(executionId); ok {
		w := waiter.(*responseWaiter) // Get the waiter struct

		// Check if request has expired
		if time.Now().After(w.expireTime) {
			r.logger.Warn().Str("execution_id", executionId).Msg("Request has expired - discarding late response")
			return nil
		}

		// Mark as handled before sending to prevent race with timeout
		if !w.expired.CompareAndSwap(false, true) {
			r.logger.Warn().Str("execution_id", executionId).Msg("Request already handled - discarding duplicate response")
			return nil
		}

		// Send the response
		w.responseChan <- RouterResponse{
			ChannelId: channelId,
			NodeData:  data,
			TimedOut:  false,
		}
		return nil
	}

	// If we get here, the request was never registered or was already cleaned up
	r.logger.Warn().Str("execution_id", executionId).Msg("No waiter found for execution ID - discarding response")
	return nil
}
