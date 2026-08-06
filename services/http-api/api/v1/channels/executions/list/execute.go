package list

import (
	"fmt"
	"time"

	gruenthttp "github.com/ongruent/gruent/libs/http"
	clickhouseService "github.com/ongruent/gruent/services/http-api/services/clickhouse"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type ChannelExecutionListItem struct {
	ExecutionId   string            `ch:"executionId" json:"executionId"`
	StartTime     time.Time         `ch:"startTime" json:"startTime"`
	EndTime       time.Time         `ch:"endTime" json:"endTime"`
	DurationMs    int64             `ch:"durationMs" json:"durationMs"`
	ChannelId     string            `ch:"channelId" json:"channelId"`
	LogAttributes map[string]string `ch:"logAttributes" json:"logAttributes"`
	Status        string            `ch:"status" json:"status"`
}

func execute(request *gruenthttp.Request[Input]) (*Output, error) {
	request.Span.AddEvent("execute")
	request.Span.SetAttributes(
		attribute.String("request.Input.ChannelId", request.Input.ChannelId),
		attribute.String("request.Input.SpaceId", request.Input.SpaceId),
		attribute.Int("request.Input.Page", request.Input.Page),
		attribute.Int("request.Input.ResultsPerPage", request.Input.ResultsPerPage),
		attribute.String("request.Input.StatusFilter", request.Input.StatusFilter),
	)

	// Query for channel executions count
	dbCtx, dbSpan := otel.Tracer("clickhouse").Start(request.Context, "Count channel executions")
	var channelExecutionsCount *int64
	countQuery := `
		SELECT toInt64(count(*))
		FROM (
			SELECT TraceId
			FROM channel_obsv.otel_logs 
			WHERE ResourceAttributes['service.instance.id'] LIKE ?
				AND SeverityText != 'DEBUG'
				AND LogAttributes['execution_id'] IS NOT NULL
				AND LogAttributes['execution_id'] != ''
				AND (? = '' OR LogAttributes['level'] = ?)
			GROUP BY TraceId)`

	err := clickhouseService.Conn.QueryRow(dbCtx, countQuery,
		request.Input.ChannelId,
		request.Input.StatusFilter,
		request.Input.StatusFilter,
	).Scan(&channelExecutionsCount)
	if err != nil {
		dbSpan.SetStatus(codes.Error, "Failed to count channel executions")
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", err.Error()),
		)
		dbSpan.End()
		return nil, err
	}

	if channelExecutionsCount == nil {
		channelExecutionsCount = new(int64)
		*channelExecutionsCount = 0
	}

	dbSpan.SetStatus(codes.Ok, "Successfully counted channel executions")
	dbSpan.SetAttributes(
		attribute.String("db.result.count", fmt.Sprintf("%d", *channelExecutionsCount)),
	)
	dbSpan.End()

	// List channel executions
	dbCtx, dbSpan = otel.Tracer("clickhouse").Start(request.Context, "List channel executions")

	channelExecutionsList := []ChannelExecutionListItem{}

	query := `
		SELECT 
			TraceId as executionId,
			min(Timestamp) as startTime,
			max(Timestamp) as endTime,
			dateDiff('millisecond', min(Timestamp), max(Timestamp)) as durationMs,
			any(ResourceAttributes['service.instance.id']) as channelId,
			any(LogAttributes) as logAttributes,
			CASE 
				WHEN MAX(LogAttributes['flag'] = 'init_error') THEN 'Init Error'
				WHEN MAX(LogAttributes['flag'] = 'timeout') THEN 'Timeout'
				WHEN MAX(LogAttributes['flag'] = 'invalid_data') THEN 'Invalid Data'
				WHEN MAX(LogAttributes['flag'] = 'complete') THEN 'Completed'
				WHEN MAX(LogAttributes['flag'] = 'fatal_error') THEN 'Fatal Error'
				ELSE 'Unknown'
			END as status
		FROM (
			SELECT *
			FROM channel_obsv.otel_logs
			WHERE TraceId IS NOT NULL
			  AND TraceId != ''
			  AND ResourceAttributes['service.instance.id'] LIKE ?
			  AND SeverityText != 'DEBUG'
			  AND LogAttributes['execution_id'] IS NOT NULL
				AND LogAttributes['execution_id'] != ''
				AND (? = '' OR LogAttributes['level'] = ?)
		) AS filtered_logs
		GROUP BY TraceId
		ORDER BY startTime DESC
		LIMIT ? OFFSET ?`

	err = clickhouseService.Conn.Select(dbCtx, &channelExecutionsList, query,
		request.Input.ChannelId,
		request.Input.StatusFilter,
		request.Input.StatusFilter,
		request.Input.ResultsPerPage,
		(request.Input.Page-1)*request.Input.ResultsPerPage,
	)
	if err != nil {
		dbSpan.SetAttributes(attribute.String("error", fmt.Sprintf("error executing ClickHouse query: %s", err)))
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	output := Output{
		ChannelExecutionsList: channelExecutionsList,
		Page:                  request.Input.Page,
		TotalPages:            (int(*channelExecutionsCount) + request.Input.ResultsPerPage - 1) / request.Input.ResultsPerPage,
		TotalResults:          int(*channelExecutionsCount),
		ResultsPerPage:        request.Input.ResultsPerPage,
	}

	return &output, nil
}
