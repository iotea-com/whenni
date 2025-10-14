package main

import (
	"context"
	"time"

	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
)

func (n *HttpSourceNode) Deinit(params node.DeinitParams) node.Error {
	// gracefully shut down the HTTP server
	if n.server != nil {
		// Create a context with timeout for shutdown
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		// Log start of shutdown
		n.logger.Debug().Msg("Starting graceful shutdown of HTTP server")

		// Create a done channel to track shutdown completion
		done := make(chan struct{})
		go func() {
			if err := n.server.ShutdownWithContext(shutdownCtx); err != nil {
				n.logger.Debug().Msgf("Failed to gracefully shutdown server: %v", err)
			}
			close(done)
		}()

		// Wait for shutdown or timeout
		select {
		case <-done:
			n.logger.Debug().Msg("HTTP server shutdown completed")
		case <-shutdownCtx.Done():
			n.logger.Debug().Msg("HTTP server shutdown timed out")
		}
	}

	// Safe channel closing function
	safeClose := func(ch chan node.IoData) {
		defer func() {
			if r := recover(); r != nil {
				n.logger.Debug().Msgf("Recovered from channel close panic: %v", r)
			}
		}()
		close(ch)
	}

	// Close all input channels
	for _, ch := range n.inputChannels {
		safeClose(ch.Channel)
	}

	// Close all output channels
	for _, ch := range n.outputChannels {
		safeClose(ch.Channel)
	}

	// Close error channel if it exists
	if n.errChan != nil {
		close(n.errChan)
	}

	return node.Error{
		Type: node.NoError,
	}
}
