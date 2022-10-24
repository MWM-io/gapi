package log

import (
	"context"
	"fmt"
	"time"

	opencensus "go.opencensus.io/trace"
	opentelemetry "go.opentelemetry.io/otel/trace"
)

// Entry represents a single log message.
// It carries all the data for your log message.
//
// Although the fields are exported, you should use the EntryOption  to modify your entries,
// that will avoid some weird override behavior.
type Entry struct {
	// Context is current context. Ie: the request context.
	// It is mainly used so that other EntryOption can fill other Entry fields. (ie: tracing)
	Context context.Context `json:"-"`

	// Timestamp is the time of the entry. If zero, the current time is used.
	Timestamp time.Time `json:"timestamp"`

	// Severity is the severity level of the log entry.
	Severity Severity `json:"severity"`

	// Message is the main message you want to log.
	Message string `json:"message"`

	// Labels are a map of <key:value> strings to contain data that could be indexed for faster queries.
	Labels map[string]string `json:"labels"`

	// Fields contains additional information to complete the message.
	// It can contain more complex data than labels
	Fields map[string]interface{} `json:"fields"`

	// TraceID is a unique identity of a trace.
	TraceID string `json:"trace_id"`
	// SpanID is a unique identity of a span in a trace.
	SpanID string `json:"span_id"`
	// IsTraceSampled indicates whether the trace is sampled or not.
	IsTraceSampled bool `json:"is_trace_sampled"`

	// StackTrace represents the stackTrace where the log was created.
	StackTrace StackTrace `json:"stack_trace"`
}

// NewEntry builds a new Entry.
func NewEntry(msg string) Entry {
	return Entry{
		Message: msg,
		Context: context.Background(),
	}
}

// EntryOption is a function that will modify an entry.
// It can be passed to a logger to be applied to all its logging calls.
// EntryOptions must check that the field they modify hasn't been set to respect the option order.
type EntryOption func(entry *Entry)

// MultiOpt combines multiple options into one EntryOption.
// Options order will be reversed so that the first one is executed last.
func MultiOpt(opts ...EntryOption) EntryOption {
	return func(entry *Entry) {
		for i := range opts {
			opts[len(opts)-i-1](entry)
		}
	}
}

// MessageOpt will set the message if it is empty.
func MessageOpt(message string) EntryOption {
	return func(entry *Entry) {
		if entry.Message == "" {
			entry.Message = message
		}
	}
}

// LabelsOpt sets labels values.
// If a label already exists, it won't be override.
func LabelsOpt(labels map[string]string) EntryOption {
	return func(entry *Entry) {
		if entry.Labels == nil {
			entry.Labels = make(map[string]string)
		}

		for key, value := range labels {
			if _, ok := entry.Labels[key]; ok {
				continue
			}

			entry.Labels[key] = value
		}
	}
}

// FieldsOpt sets fields values.
// If a field already exists, it won't be override.
func FieldsOpt(fields map[string]interface{}) EntryOption {
	return func(entry *Entry) {
		if entry.Fields == nil {
			entry.Fields = make(map[string]interface{})
		}

		for key, value := range fields {
			if _, ok := entry.Fields[key]; ok {
				continue
			}

			entry.Fields[key] = value
		}
	}
}

// SeverityOpt sets the severity if it is still DefaultSeverity.
func SeverityOpt(severity Severity) EntryOption {
	return func(entry *Entry) {
		if entry.Severity == DefaultSeverity {
			entry.Severity = severity
		}
	}
}

// ContextOpt sets the context if it is still context.Background().
func ContextOpt(ctx context.Context) EntryOption {
	return func(entry *Entry) {
		if entry.Context == context.Background() {
			entry.Context = ctx
		}
	}
}

// TimestampOpt sets the timestamp it is still time.IsZero().
func TimestampOpt(t time.Time) EntryOption {
	return func(entry *Entry) {
		if entry.Timestamp.IsZero() {
			entry.Timestamp = t
		}
	}
}

// StackTraceOpt sets the stacktrace if it is still nil.
func StackTraceOpt(trace StackTrace) EntryOption {
	return func(entry *Entry) {
		if entry.StackTrace == nil {
			entry.StackTrace = trace
		}
	}
}

// TracingOpt sets the tracing data it is not yet set.
func TracingOpt(traceID, spanID string, isSampled bool) EntryOption {
	return func(entry *Entry) {
		if entry.TraceID != "" || entry.SpanID != "" {
			return
		}

		entry.TraceID = traceID
		entry.SpanID = spanID
		entry.IsTraceSampled = isSampled
	}
}

// TimestampNowOpt sets the timestamp of an entry to time.Now()
func TimestampNowOpt() EntryOption {
	return TimestampOpt(time.Now())
}

// OpencensusTraceOpt sets the trace information from the opencensus context.
func OpencensusTraceOpt() EntryOption {
	return func(entry *Entry) {
		spanContext := opencensus.FromContext(entry.Context).SpanContext()
		traceID := spanContext.TraceID.String()
		spanID := spanContext.SpanID.String()

		if traceID == "00000000000000000000000000000000" || spanID == "0000000000000000" {
			return
		}

		TracingOpt(traceID, spanID, spanContext.IsSampled())(entry)
	}
}

// OpentelemetryTraceOpt sets the trace information from the opentelemetry context.
func OpentelemetryTraceOpt() EntryOption {
	return func(entry *Entry) {
		spanContext := opentelemetry.SpanContextFromContext(entry.Context)
		traceID := spanContext.TraceID().String()
		spanID := spanContext.SpanID().String()

		if traceID == "00000000000000000000000000000000" || spanID == "0000000000000000" {
			return
		}

		TracingOpt(traceID, spanID, spanContext.IsSampled())(entry)
	}
}

// AnyOpt will try to cast the given argument to extract useful data
// It uses single-purpose interfaces to determine its fields:
// - error, use error.Error() as message and set severity to error (if not set)
// - interface{ Context() context.Context }
// - fmt.Stringer for message.
// - interface{ Message() string }
// - interface{ Timestamp() time.Time }
// - interface{ Severity() Severity }
// - interface{ Labels() map[string]string }
// - interface{ Fields() map[string]interface{} }
// - interface {
//		TraceID() string
//		SpanID() string
//		IsTraceSampled() bool
//	}
// - interface{ StackTrace() StackTrace }
func AnyOpt(v interface{}) EntryOption {
	return func(entry *Entry) {
		var opts []EntryOption

		if entryErr, ok := v.(error); ok {
			opts = append(opts, SeverityOpt(ErrorSeverity))
			opts = append(opts, MessageOpt(entryErr.Error()))
		}

		withMessage, ok := v.(interface{ Message() string })
		if ok {
			opts = append(opts, MessageOpt(withMessage.Message()))
		} else if entry.Message == "" {
			if stringer, ok := v.(fmt.Stringer); ok {
				opts = append(opts, MessageOpt(stringer.String()))
			} else {
				opts = append(opts, MessageOpt(fmt.Sprintf("%v", v)))
			}
		}

		withContext, ok := v.(interface{ Context() context.Context })
		if ok {
			opts = append(opts, ContextOpt(withContext.Context()))
		}

		withTimestamp, ok := v.(interface{ Timestamp() time.Time })
		if ok {
			opts = append(opts, TimestampOpt(withTimestamp.Timestamp()))
		}

		withSeverity, ok := v.(interface{ Severity() Severity })
		if ok {
			opts = append(opts, SeverityOpt(withSeverity.Severity()))
		}

		withLabels, ok := v.(interface{ Labels() map[string]string })
		if ok {
			opts = append(opts, LabelsOpt(withLabels.Labels()))
		}

		withFields, ok := v.(interface{ Fields() map[string]interface{} })
		if ok {
			opts = append(opts, FieldsOpt(withFields.Fields()))
		}

		withTracing, ok := v.(interface {
			TraceID() string
			SpanID() string
			IsTraceSampled() bool
		})
		if ok {
			opts = append(opts, TracingOpt(withTracing.TraceID(), withTracing.SpanID(), withTracing.IsTraceSampled()))
		}

		withStack, ok := v.(interface{ StackTrace() StackTrace })
		if ok {
			opts = append(opts, StackTraceOpt(withStack.StackTrace()))
		}

		MultiOpt(opts...)(entry)
	}
}
