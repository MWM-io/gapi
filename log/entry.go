package log

import (
	"context"
	"fmt"
	"runtime"
	"time"

	opencensus "go.opencensus.io/trace"
	opentelemetry "go.opentelemetry.io/otel/trace"

	"github.com/mwm-io/gapi/stacktrace"
)

type StackTrace interface {
	fmt.Stringer
	Frames() []runtime.Frame
	Last() (runtime.Frame, bool)
}

// Entry contains all the data of a log entry.
type Entry struct {
	// Context is current context. Ie: the request context.
	Context context.Context

	// Timestamp is the time of the entry. If zero, the current time is used.
	Timestamp time.Time

	// Severity is the severity level of the log entry.
	Severity Severity

	// Message is the main message you want to log.
	// For easier filtering/querying, use generic messages without variables, and use fields to store the variables.
	Message string

	// Labels are a map of <key:value> strings to contain data that could be indexed for faster queries.
	Labels map[string]string
	// Fields contains additional information to complete the message.
	// It can contain more complex data than labels
	Fields map[string]interface{}

	// TraceID is a unique identity of a trace.
	TraceID string
	// SpanID is a unique identity of a span in a trace.
	SpanID string
	// IsTraceSampled indicates whether the trace is sampled or not.
	IsTraceSampled bool
	// StackTrace represents the stackTrace where the log was created.
	StackTrace StackTrace
}

func NewEntry(msg string) Entry {
	return Entry{
		Message: msg,
		Context: context.Background(),
	}
}

// EntryOption is a function that will add information to an entry.
// EntryOptions must check that the field they modify hasn't been set to respect the option order.
type EntryOption func(entry *Entry)

// MultiOpt combine multiple options into one EntryOption.
// Options order will be inversed so that the first one is executed first.
func MultiOpt(opts ...EntryOption) EntryOption {
	return func(entry *Entry) {
		for i := range opts {
			opts[len(opts)-i-1](entry)
		}
	}
}

func MessageOpt(message string) EntryOption {
	return func(entry *Entry) {
		if entry.Message == "" {
			entry.Message = message
		}
	}
}

// LabelsOpt set default labels.
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

// FieldsOpt set default fields.
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

// SeverityOpt set the default Severity.
func SeverityOpt(severity Severity) EntryOption {
	return func(entry *Entry) {
		if entry.Severity == DefaultSeverity {
			entry.Severity = severity
		}
	}
}

func ContextOpt(ctx context.Context) EntryOption {
	return func(entry *Entry) {
		if entry.Context == context.Background() {
			entry.Context = ctx
		}
	}
}

func TimestampOpt(t time.Time) EntryOption {
	return func(entry *Entry) {
		if entry.Timestamp.IsZero() {
			entry.Timestamp = t
		}
	}
}

func StackTraceOpt(trace StackTrace) EntryOption {
	return func(entry *Entry) {
		if entry.StackTrace == nil {
			entry.StackTrace = trace
		}
	}
}

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

// TimestampNowOpt set the default timestamp of an entry to time.Now() when applied.
// (ie: when the entry will be built)
func TimestampNowOpt() EntryOption {
	return TimestampOpt(time.Now())
}

// DefaultStackTraceOpt set the default stack trace of an entry.
func DefaultStackTraceOpt() EntryOption {
	return StackTraceOpt(stacktrace.New())
}

// OpencensusTraceOpt set the trace information from the opencensus context.
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

// OpentelemetryTraceOpt set the trace information from the opentelemetry context.
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
// - interface{ TraceID() string }
// - interface{ SpanID() string }
// - interface{ IsTraceSampled() bool }
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
