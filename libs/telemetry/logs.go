package telemetry

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/iotea-com/iotea/libs/legacy/engine/environment"
	"github.com/iotea-com/iotea/services/http-api/config"
	"go.opentelemetry.io/contrib/bridges/otelzap"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

// Global logger instance. Should be used for logs that are not associated with a specific
// request/trace, require auditing through text search, or that contain a large payload.
var Logger *zap.Logger

// Sets up the logger provider. Logs will always be written to stdout via zap. If an otel collector
// connection has been made, logs will also be exported to the remote collector.
func setupLoggerProvider(ctx context.Context, conn *grpc.ClientConn, res *resource.Resource) (*sdklog.LoggerProvider, error) {
	// Setup remote logging if possible
	var remoteExporter sdklog.Exporter
	if conn != nil {
		re, err := otlploggrpc.New(ctx, otlploggrpc.WithGRPCConn(conn))
		if err != nil {
			return nil, err
		}
		remoteExporter = re
	}

	// Create the LoggerProvider with the exporters and resource
	opts := []sdklog.LoggerProviderOption{sdklog.WithResource(res)}
	if remoteExporter != nil {
		opts = append(opts, sdklog.WithProcessor(sdklog.NewBatchProcessor(remoteExporter)))
	}
	lp := sdklog.NewLoggerProvider(opts...)
	return lp, nil
}

// Creates a new zap logger. If an otel logger provider is provided, it will be used to create the
// logger.
func NewLogger(loggerProvider *sdklog.LoggerProvider) (*zap.Logger, error) {
	logger, err := (func() (*zap.Logger, error) {
		if config.Env == environment.Production {
			return zap.NewProduction()
		}

		loggerConfig := zap.NewDevelopmentConfig()
		loggerConfig.DisableCaller = true
		loggerConfig.DisableStacktrace = true
		loggerConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		loggerConfig.Encoding = "console"
		loggerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		return loggerConfig.Build()
	})()
	if err != nil {
		return nil, err
	}

	if loggerProvider != nil {
		otelCore := otelzap.NewCore("http-api", otelzap.WithLoggerProvider(loggerProvider))
		logger = zap.New(zapcore.NewTee(logger.Core(), otelCore))
	}

	return logger, nil
}

// Returns a logger from the request context. This will automatically associate trace/span IDs
// with logs.
func GetLoggerFromContext(c *fiber.Ctx) *zap.Logger {
	return Logger.With(
		zap.String("path", c.Path()),
		zap.String("method", c.Method()),
		zap.Any("otel_context", c.UserContext()),
	)
}
