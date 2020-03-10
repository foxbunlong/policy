package mysql

import (
	"context"
	"runtime"

	"github.com/opentracing/opentracing-go"
)

func startSpan(ctx context.Context, name, sql string) (s opentracing.Span) {

	if span := opentracing.SpanFromContext(ctx); span != nil {
		c, f, l, _ := runtime.Caller(1)
		tracer := opentracing.GlobalTracer()
		s = tracer.StartSpan(name, opentracing.ChildOf(span.Context()))
		s.SetTag("component", "mariadb")
		s.SetTag("sql", sql)
		s.SetTag("func.c", c)

		s.SetTag("func.f", f)

		s.SetTag("func.l", l)

		carier := opentracing.TextMapCarrier{"sql": sql}
		err := tracer.Inject(s.Context(), opentracing.TextMap, carier)
		if err != nil {

		}
	}
	return
}

func finishSpan(span opentracing.Span) {
	if span != nil {
		span.Finish()
	}
}
