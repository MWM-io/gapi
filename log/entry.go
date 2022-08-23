package log

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"time"

	opencensus "go.opencensus.io/trace"
	opentelemetry "go.opentelemetry.io/otel/trace"
)

// Entry contains all the data of a log entry.
type Entry struct {
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
	// StackTrace is a slice of program counters. (ie: built by runtime.Callers)
	StackTrace []uintptr
}

// NewEntry tries to build a new entry from interface{}thing.
// It uses single-purpose interfaces to determine its fields:
// - interface{ Message() string }
// - interface{ Timestamp() time.Time }
// - interface{ Severity() Severity }
// - interface{ Labels() map[string]string }
// - interface{ Fields() map[string]interface{} }
// - interface{ TraceID() string }
// - interface{ SpanID() string }
// - interface{ IsTraceSampled() bool }
// - interface{ Stack() []uintptr }
func NewEntry(v interface{}) Entry {
	entry := Entry{}

	withMessage, ok := v.(interface{ Message() string })
	if ok {
		entry.Message = withMessage.Message()
	} else {
		if stringer, ok := v.(fmt.Stringer); ok {
			entry.Message = stringer.String()
		} else {
			entry.Message = fmt.Sprintf("%v", v)
		}
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

	withStack, ok := v.(interface{ Stack() []uintptr })
	if ok {
		entry.StackTrace = withStack.Stack()
	}

	return entry
}

// StackInfo contains basic data to describe a runtime.Frame
type StackInfo struct {
	File     string
	Function string
	Line     int
}

// GetLastStackInfo return stack information on the last function called before getting inside the log package.
func (e Entry) GetLastStackInfo() StackInfo {
	if e.StackTrace == nil {
		return StackInfo{}
	}

	frames := runtime.CallersFrames(e.StackTrace[:])
	for {
		frame, more := frames.Next()

		if strings.Contains(frame.File, "gapi/log") && more {
			continue
		}

		return StackInfo{
			File:     frame.File,
			Function: frame.Function,
			Line:     frame.Line,
		}
	}
}

// GetStackInfo returns the full stack info.
func (e Entry) GetStackInfo() []StackInfo {
	var stack []StackInfo

	frames := runtime.CallersFrames(e.StackTrace[:])
	for {
		frame, more := frames.Next()

		stack = append(stack, StackInfo{
			File:     frame.File,
			Function: frame.Function,
			Line:     frame.Line,
		})

		if !more {
			break
		}
	}

	return stack
}

// EntryOption is a function that will add information to an entry.
// It takes the current context where the entry was created.
// If no context was passed during the log creation, it will be an empty context.Background()
// An EntryOption shouldn't override the entry existing fields, it should only add default values.
type EntryOption func(ctx context.Context, entry *Entry)

// LabelsOpt set default labels.
// If this option was already applied with labels with the same key, they will take the last value.
func LabelsOpt(labels map[string]string) EntryOption {
	return func(ctx context.Context, entry *Entry) {
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
// If this option was already applied with fields with the same key, they will take the last value.
func FieldsOpt(fields map[string]interface{}) EntryOption {
	return func(ctx context.Context, entry *Entry) {
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
	return func(ctx context.Context, entry *Entry) {
		if entry.Severity != DefaultSeverity {
			return
		}

		entry.Severity = severity
	}
}

// TimestampNowOpt set the default timestamp of an entry to time.Now() when applied.
// (ie: when the entry will be built)
func TimestampNowOpt() EntryOption {
	return func(ctx context.Context, entry *Entry) {
		if !entry.Timestamp.IsZero() {
			return
		}

		entry.Timestamp = time.Now()
	}
}

// DefaultStackTraceOpt set the default stack trace of an entry.
func DefaultStackTraceOpt() EntryOption {
	return func(ctx context.Context, entry *Entry) {
		if entry.StackTrace != nil {
			return
		}

		pc := make([]uintptr, 32)
		n := runtime.Callers(2, pc)
		entry.StackTrace = pc[:n]
	}
}

// OpencensusTraceOpt set the trace information from the opencensus context.
func OpencensusTraceOpt() EntryOption {
	return func(ctx context.Context, entry *Entry) {
		if entry.TraceID != "" || entry.SpanID != "" {
			return
		}

		spanContext := opencensus.FromContext(ctx).SpanContext()
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
	return func(ctx context.Context, entry *Entry) {
		if entry.TraceID != "" || entry.SpanID != "" {
			return
		}

		spanContext := opentelemetry.SpanContextFromContext(ctx)
		if traceID := spanContext.TraceID().String(); traceID != "00000000000000000000000000000000" {
			entry.TraceID = traceID
		}
		if spanID := spanContext.SpanID().String(); spanID != "0000000000000000" {
			entry.SpanID = spanID
		}
		entry.IsTraceSampled = spanContext.IsSampled()
	}
}
