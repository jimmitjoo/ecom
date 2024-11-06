package tracing

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// mockTracer skapar en tracer som skriver till minnet istället för Jaeger
func mockTracer(t *testing.T) (*sdktrace.TracerProvider, error) {
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
	)

	otel.SetTracerProvider(tp)
	return tp, nil
}

func TestInitTracer(t *testing.T) {
	serviceName := "test-service"
	ctx := context.Background()

	tp, err := mockTracer(t)
	assert.NoError(t, err)
	assert.NotNil(t, tp)
	defer func() {
		err := tp.Shutdown(ctx)
		assert.NoError(t, err)
	}()

	globalTP := otel.GetTracerProvider()
	assert.NotNil(t, globalTP)

	tracer := otel.Tracer(serviceName)
	assert.NotNil(t, tracer)

	_, span := tracer.Start(ctx, "test-span")
	assert.NotNil(t, span)
	defer span.End()

	spanContext := span.SpanContext()
	assert.True(t, spanContext.IsValid())
}

func TestTracerAttributes(t *testing.T) {
	serviceName := "test-service"
	ctx := context.Background()

	tp, err := mockTracer(t)
	assert.NoError(t, err)
	defer func() {
		err := tp.Shutdown(ctx)
		assert.NoError(t, err)
	}()

	tracer := otel.Tracer(serviceName)
	_, span := tracer.Start(ctx, "test-span-with-attributes")
	defer span.End()

	span.SetAttributes(
		attribute.String("string-key", "string-value"),
		attribute.Int("int-key", 123),
		attribute.Bool("bool-key", true),
	)

	assert.True(t, span.IsRecording())
}

func TestTracerWithEvents(t *testing.T) {
	serviceName := "test-service"
	ctx := context.Background()

	tp, err := mockTracer(t)
	assert.NoError(t, err)
	defer func() {
		err := tp.Shutdown(ctx)
		assert.NoError(t, err)
	}()

	tracer := otel.Tracer(serviceName)
	_, span := tracer.Start(ctx, "test-span-with-events")
	defer span.End()

	// Lägg till events med attribut
	span.AddEvent("event-1")
	span.AddEvent("event-2", trace.WithAttributes(
		attribute.String("key1", "value1"),
		attribute.Int("key2", 123),
	))

	assert.True(t, span.IsRecording())
}

func TestTracerSpanHierarchy(t *testing.T) {
	serviceName := "test-service"
	ctx := context.Background()

	tp, err := mockTracer(t)
	assert.NoError(t, err)
	defer tp.Shutdown(ctx)

	tracer := otel.Tracer(serviceName)

	ctx, parentSpan := tracer.Start(ctx, "parent-span")
	defer parentSpan.End()

	_, childSpan := tracer.Start(ctx, "child-span")
	defer childSpan.End()

	parentContext := parentSpan.SpanContext()
	childContext := childSpan.SpanContext()

	assert.NotEqual(t, parentContext.SpanID(), childContext.SpanID())
	assert.Equal(t, parentContext.TraceID(), childContext.TraceID())
}
