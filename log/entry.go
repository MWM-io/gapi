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
type EntryOption func(entry *Entry)

// LabelsOpt set default labels.
func LabelsOpt(labels map[string]string) EntryOption {
	return func(entry *Entry) {
		if entry.Labels == nil {
			entry.Labels = make(map[string]string)
		}

		for key, value := range labels {
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
			entry.Fields[key] = value
		}
	}
}

// SeverityOpt set the default Severity.
func SeverityOpt(severity Severity) EntryOption {
	return func(entry *Entry) {
		entry.Severity = severity
	}
}

func ContextOpt(ctx context.Context) EntryOption {
	return func(entry *Entry) {
		entry.Context = ctx
	}
}

// TimestampNowOpt set the default timestamp of an entry to time.Now() when applied.
// (ie: when the entry will be built)
func TimestampNowOpt() EntryOption {
	return func(entry *Entry) {
		entry.Timestamp = time.Now()
	}
}

// DefaultStackTraceOpt set the default stack trace of an entry.
func DefaultStackTraceOpt() EntryOption {
	return func(entry *Entry) {
		if entry.StackTrace != nil {
			return
		}

		entry.StackTrace = stacktrace.New()
	}
}

// OpencensusTraceOpt set the trace information from the opencensus context.
func OpencensusTraceOpt() EntryOption {
	return func(entry *Entry) {
		if entry.TraceID != "" || entry.SpanID != "" {
			return
		}

		spanContext := opencensus.FromContext(entry.Context).SpanContext()
		if traceID := spanContext.TraceID.String(); traceID != "00000000000000000000000000000000" {
			entry.TraceID = traceID
		}
		if spanID := spanContext.SpanID.String(); spanID != "0000000000000000" {
			entry.SpanID = spanID
		}
		entry.IsTraceSampled = spanContext.IsSampled()

	}
}

// OpentelemetryTraceOpt set the trace information from the opentelemetry context.
func OpentelemetryTraceOpt() EntryOption {
	return func(entry *Entry) {
		if entry.TraceID != "" || entry.SpanID != "" {
			return
		}

		spanContext := opentelemetry.SpanContextFromContext(entry.Context)
		if traceID := spanContext.TraceID().String(); traceID != "00000000000000000000000000000000" {
			entry.TraceID = traceID
		}
		if spanID := spanContext.SpanID().String(); spanID != "0000000000000000" {
			entry.SpanID = spanID
		}
		entry.IsTraceSampled = spanContext.IsSampled()
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
		if entryErr, ok := v.(error); ok {
			if entry.Severity == DefaultSeverity {
				entry.Severity = ErrorSeverity
			}
			entry.Message = entryErr.Error()
		}

		withMessage, ok := v.(interface{ Message() string })
		if ok {
			entry.Message = withMessage.Message()
		} else if entry.Message == "" {
			if stringer, ok := v.(fmt.Stringer); ok {
				entry.Message = stringer.String()
			} else {
				entry.Message = fmt.Sprintf("%v", v)
			}
		}

		withContext, ok := v.(interface{ Context() context.Context })
		if ok {
			entry.Context = withContext.Context()
		}

		withTimestamp, ok := v.(interface{ Timestamp() time.Time })
		if ok {
			entry.Timestamp = withTimestamp.Timestamp()
		}

		withSeverity, ok := v.(interface{ Severity() Severity })
		if ok {
			entry.Severity = withSeverity.Severity()
		}

		withLabels, ok := v.(interface{ Labels() map[string]string })
		if ok {
			entry.Labels = withLabels.Labels()
		}

		withFields, ok := v.(interface{ Fields() map[string]interface{} })
		if ok {
			entry.Fields = withFields.Fields()
		}

		withTraceID, ok := v.(interface{ TraceID() string })
		if ok {
			entry.TraceID = withTraceID.TraceID()
		}

		withSpanID, ok := v.(interface{ SpanID() string })
		if ok {
			entry.SpanID = withSpanID.SpanID()
		}

		withIsTraceSampled, ok := v.(interface{ IsTraceSampled() bool })
		if ok {
			entry.IsTraceSampled = withIsTraceSampled.IsTraceSampled()
		}

		withStack, ok := v.(interface{ StackTrace() StackTrace })
		if ok {
			entry.StackTrace = withStack.StackTrace()
		}
	}
}
