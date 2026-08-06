package get

import (
	"fmt"
	"time"

	gruenthttp "github.com/ongruent/gruent/libs/http"
	clickhouseService "github.com/ongruent/gruent/services/http-api/services/clickhouse"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type NodeExecutionLog struct {
	TraceId            string            `ch:"traceId" json:"traceId"`
	SpanId             string            `ch:"spanId" json:"spanId"`
	Timestamp          time.Time         `ch:"timestamp" json:"timestamp"`
	Level              string            `ch:"level" json:"level"`
	Body               string            `ch:"body" json:"body"`
	LogAttributes      map[string]string `ch:"logAttributes" json:"logAttributes"`
	ResourceAttributes map[string]string `ch:"resourceAttributes" json:"resourceAttributes"`
	Status             string            `ch:"status" json:"status"`
}

type ChannelExecution struct {
	TraceId           string             `ch:"executionId" json:"executionId"`
	StartTime         time.Time          `ch:"startTime" json:"startTime"`
	EndTime           time.Time          `ch:"endTime" json:"endTime"`
	DurationMs        int64              `ch:"durationMs" json:"durationMs"`
	ChannelId         string             `ch:"channelId" json:"channelId"`
	Status            string             `ch:"status" json:"status"`
	NodeExecutionLogs []NodeExecutionLog `ch:"nodeExecutionLogs" json:"nodeExecutionLogs"`
}

func execute(request *gruenthttp.Request[Input]) (*Output, error) {
	request.Span.AddEvent("execute")
	request.Span.SetAttributes(
		attribute.String("request.Input.ChannelExecutionId", request.Input.ChannelExecutionId),
		attribute.String("request.Input.SpaceId", request.Input.SpaceId),
	)

	// Get node execution logs in a channel execution
	dbCtx, dbSpan := otel.Tracer("clickhouse").Start(request.Context, "List node execution logs")

	nodeExecutionLogs := []NodeExecutionLog{}

	query := `
	SELECT
		TraceId as traceId,
		SpanId as spanId,
		Timestamp as timestamp,
		SeverityText as level,
		Body as body,
		LogAttributes as logAttributes,
		ResourceAttributes as resourceAttributes,
		CASE 
			WHEN LogAttributes['flag'] = 'init_error' THEN 'Init Error'
			WHEN LogAttributes['flag'] = 'timeout' THEN 'Timeout'
			WHEN LogAttributes['flag'] = 'invalid_data' THEN 'Invalid Data'
			WHEN LogAttributes['flag'] = 'complete' THEN 'Completed'
			ELSE 'Unknown'
		END as status
	FROM channel_obsv.otel_logs 
	WHERE TraceId = ?
	ORDER BY Timestamp DESC
	`

	err := clickhouseService.Conn.Select(dbCtx, &nodeExecutionLogs, query, request.Input.ChannelExecutionId)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error executing ClickHouse query: %s", err)),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	// Group logs into a channel execution
	durationMs := int64(nodeExecutionLogs[len(nodeExecutionLogs)-1].Timestamp.Sub(nodeExecutionLogs[0].Timestamp).Milliseconds())
	channelId := nodeExecutionLogs[0].ResourceAttributes["service.instance.id"]
	channelStatus := (func() string {
		for _, log := range nodeExecutionLogs {
			if log.Status != "Unknown" {
				return log.Status
			}
		}
		return "Unknown"
	})()
	channelExecution := &ChannelExecution{
		TraceId:           request.Input.ChannelExecutionId,
		ChannelId:         channelId,
		StartTime:         nodeExecutionLogs[0].Timestamp,
		EndTime:           nodeExecutionLogs[len(nodeExecutionLogs)-1].Timestamp,
		DurationMs:        durationMs,
		Status:            channelStatus,
		NodeExecutionLogs: nodeExecutionLogs,
	}

	output := Output{
		ChannelExecution: channelExecution,
	}

	return &output, nil
}
