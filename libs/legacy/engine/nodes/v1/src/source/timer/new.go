package main

import (
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	timerSourceNodeConfig "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/source/timer/config"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

const (
	NodeName = "Timer Trigger"
)

type TimerSourceNode struct {
	config         timerSourceNodeConfig.TimerSourceNodeConfig
	inputChannels  []node.IoChannel
	outputChannels []node.IoChannel
	notifyChannel  chan node.Notification

	// CRON scheduler
	triggerChannel chan bool
	cronScheduler  *cron.Cron

	// Observability
	tracer trace.Tracer
	meter  metric.Meter
	logger zerolog.Logger
}

func New() node.Interface {
	return &TimerSourceNode{
		inputChannels: []node.IoChannel{
			{
				Id:      "passThrough",
				Channel: make(chan node.IoData, 100),
				Type:    []node.DataType{node.SignalDataType},
			},
		},
		outputChannels: []node.IoChannel{
			{
				Id:      "trigger",
				Channel: make(chan node.IoData, 100),
				Type:    []node.DataType{node.SignalDataType},
			},
		},
	}
}

func main() {} // main is required by go, but node is built as plugin
