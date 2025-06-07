package instrumentation

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/vnnyx/employee-management/internal/auth/entity"
	"github.com/vnnyx/employee-management/internal/constants"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// tracer is the OpenTelemetry tracer instance used for creating spans.
var tracer = otel.Tracer("github.com/vnnyx/employee-management")

func NewTraceSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	if info, ok := ctx.Value(constants.KeyFiberCtxInformation).(entity.FiberCtxInformation); ok && info.Enable {
		var timestamp atomic.Value
		timestamp.Store(time.Now().In(time.Local).Format(`15:04:05`))

		fmt.Printf("[%s] %s %s \"%s\"\n", timestamp.Load().(string), info.Method, info.OriginalURL, name)
	}

	return tracer.Start(
		ctx,
		name,
	)
}

// NewTraceSpanWithBaggage starts a new OpenTelemetry trace span with the specified name and
// attaches baggage from the context as span attributes.
//
// Parameters:
//
//	ctx  - The context containing baggage and optional trace metadata.
//	name - The name of the span.
//
// Returns:
//
//	A new context containing the created span and the span itself.
func NewTraceSpanWithBaggage(ctx context.Context, name string) (context.Context, trace.Span) {
	return tracer.Start(
		ctx,
		name,
		trace.WithAttributes(baggageToAttributes(ctx)...),
	)
}

// RecordSpanError records an error on the provided span and sets its status to Error.
//
// Parameters:
//
//	span - The span on which to record the error.
//	err  - The error to record.
func RecordSpanError(span trace.Span, err error) {
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}

// RecordSpanErrorWithStackTrace records an error on the provided span with a stack trace
// and sets its status to Error.
//
// Parameters:
//
//	span - The span on which to record the error.
//	err  - The error to record.
func RecordSpanErrorWithStackTrace(span trace.Span, err error) {
	span.RecordError(err, trace.WithStackTrace(true))
	span.SetStatus(codes.Error, err.Error())
}

// baggageToAttributes extracts baggage from the context and converts each baggage member into
// an OpenTelemetry attribute. This allows baggage values to be attached to spans.
//
// Parameters:
//
//	ctx - The context containing baggage.
//
// Returns:
//
//	A slice of attributes representing the baggage.
func baggageToAttributes(ctx context.Context) []attribute.KeyValue {
	bag := baggage.FromContext(ctx)
	members := bag.Members()
	attributes := make([]attribute.KeyValue, 0, len(members))
	for _, member := range members {
		attributes = append(attributes, attribute.String(member.Key(), member.Value()))
	}
	return attributes
}
