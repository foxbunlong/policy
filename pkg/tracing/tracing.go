package tracing

import (
	"context"

	opentracing "github.com/opentracing/opentracing-go"
)

// SpanContext workaround for https://github.com/opentracing/specification/blob/master/rfc/trace_identifiers.md#specification-changes
type SpanContext interface {
	opentracing.SpanContext
	ToTraceID() string
	ToSpanID() string
}

// Tracer workaround for https://github.com/opentracing/specification/blob/master/rfc/trace_identifiers.md#specification-changes
type Tracer interface {
	opentracing.Tracer
	Identify(opentracing.SpanContext) SpanContext
}

// Trace contains trace id and span id for integration with other system such as log
type Trace struct {
	ID     string `json:"id"`
	SpanID string `json:"span_id"`
}

// FromContext try to extract Trace from context.Context
func FromContext(ctx context.Context) (t Trace) {
	t = Trace{}

	tracer, ok := opentracing.GlobalTracer().(Tracer)
	if !ok {
		return
	}

	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		return
	}

	spanCtx := tracer.Identify(span.Context())
	if spanCtx == nil {
		return
	}

	t.ID = spanCtx.ToTraceID()
	t.SpanID = spanCtx.ToSpanID()
	return
}
