package telemetry

import (
	"context"
	"fmt"
	"strings"

	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

// Sets up the tracer provider. If no otel collector connection is provided, a noop tracer provider
// will be used.
func setupTracerProvider(ctx context.Context, conn *grpc.ClientConn, res *resource.Resource) (*sdkTrace.TracerProvider, error) {
	if conn == nil {
		tp := sdkTrace.NewTracerProvider(
			sdkTrace.WithSampler(sdkTrace.NeverSample()),
			sdkTrace.WithResource(res),
		)
		return tp, nil
	}

	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	tp := sdkTrace.NewTracerProvider(
		sdkTrace.WithSampler(sdkTrace.AlwaysSample()),
		sdkTrace.WithResource(res),
		sdkTrace.WithSpanProcessor(sdkTrace.NewBatchSpanProcessor(traceExporter)),
	)

	return tp, nil
}

// Creates a middleware that automatically sets up tracing for the request. Skips routes that don't
// need tracing.
func TracerMiddleware() fiber.Handler {
	return otelfiber.Middleware(
		otelfiber.WithNext(func(c *fiber.Ctx) bool {
			// Skip metrics for health checks and other noisy endpoints
			return strings.HasPrefix(c.Path(), "/healthcheck")
		}),
	)
}

// Returns the span from the gin request context. This can always be treated as a parent span for
// starting new spans with `StartSpan`.
func GetRequestSpan(fiberCtx *fiber.Ctx) (context.Context, trace.Span) {
	ctx := fiberCtx.UserContext()
	return ctx, trace.SpanFromContext(ctx)
}

// Starts a new child span with the given name and returns the context and span. If the parent
// context is nil, a new background context will be created and a new span will be started.
func StartSpan(parentCtx context.Context, name string) (context.Context, trace.Span) {
	if parentCtx == nil {
		childCtx := context.Background()
		return childCtx, trace.SpanFromContext(childCtx)
	}

	parentSpan := trace.SpanFromContext(parentCtx)
	provider := parentSpan.TracerProvider()
	return provider.Tracer("http-api").Start(parentCtx, name)
}
