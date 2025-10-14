package main

import (
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
)

func (n *TimerSourceNode) Deinit(params node.DeinitParams) node.Error {
	_, deinitSpan := n.tracer.Start(params.Ctx, "Deinitialize Node")
	defer deinitSpan.End()

	// Close nodes safely
	safeClose := func(ch chan node.IoData) {
		defer func() {
			recover()
		}()
		close(ch)
	}

	// Close CRON scheduler
	if n.cronScheduler != nil {
		schedulerCtx := n.cronScheduler.Stop()
		<-schedulerCtx.Done()
	}

	// Close all input nodes
	for _, ch := range n.inputChannels {
		safeClose(ch.Channel)
	}

	// Close all output nodes
	for _, ch := range n.outputChannels {
		safeClose(ch.Channel)
	}

	return node.Error{
		Type: node.NoError,
	}
}
