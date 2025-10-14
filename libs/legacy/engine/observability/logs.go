package observability

import (
	"context"
	"fmt"

	"github.com/agoda-com/opentelemetry-logs-go/exporters/otlp/otlplogs"
	"github.com/agoda-com/opentelemetry-logs-go/exporters/otlp/otlplogs/otlplogsgrpc"
	"github.com/agoda-com/opentelemetry-logs-go/sdk/logs"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/sdk/resource"
	"google.golang.org/grpc"
)

func setupLoggerProvider(ctx context.Context, conn *grpc.ClientConn, res *resource.Resource) (*logs.LoggerProvider, error) {
	// Initialize the OTLP logs exporter using the provided gRPC connection
	client := otlplogsgrpc.NewClient(otlplogsgrpc.WithGRPCConn(conn))
	logExporter, err := otlplogs.NewExporter(ctx, otlplogs.WithClient(client))
	if err != nil {
		return nil, fmt.Errorf("failed to create logs exporter: %w", err)
	}

	// Create the LoggerProvider with the exporter and resource
	lp := logs.NewLoggerProvider(
		logs.WithBatcher(logExporter),
		logs.WithResource(res),
	)

	return lp, nil
}

type LogLevel string

const (
	LogLevelTrace LogLevel = "trace"
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

func (l LogLevel) String() string {
	return string(l)
}

func (l LogLevel) Level() zerolog.Level {
	switch l {
	case LogLevelTrace:
		return zerolog.TraceLevel
	case LogLevelDebug:
		return zerolog.DebugLevel
	case LogLevelInfo:
		return zerolog.InfoLevel
	case LogLevelWarn:
		return zerolog.WarnLevel
	case LogLevelError:
		return zerolog.ErrorLevel
	default:
		log.Warn().Msgf("Invalid log level: %s", l)
		return zerolog.InfoLevel
	}
}
