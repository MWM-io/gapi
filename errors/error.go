package errors

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"time"

	gLog "github.com/mwm-io/gapi/log"
	"github.com/mwm-io/gapi/stacktrace"
)

// Error represents the interface for the FullError struct.
// It is necessary so we can compare nil errors. (nil.(*FullError) != nil)
type Error interface {
	error
	Unwrap() error
	json.Marshaler
	xml.Marshaler
	Message() string
	Kind() string
	StatusCode() int
	Severity() gLog.Severity
	Timestamp() time.Time
	StackTrace() gLog.StackTrace

	WithMessage(string) Error
	WithKind(string) Error
	WithStatus(int) Error
	WithSeverity(gLog.Severity) Error
	WithTimestamp(time.Time) Error
	WithStackTrace(trace gLog.StackTrace) Error
}

// FullError is a concrete error that implements the Error interface
type FullError struct {
	userMessage  string
	kind         string
	errorMessage string
	status       int
	severity     gLog.Severity
	timestamp    time.Time
	stackTrace   gLog.StackTrace
	prev         error
}

// Wrap will wrap the given error and return a new Error.
func Wrap(previousError error, message string) Error {
	if previousError == nil {
		return nil
	}

	err := &FullError{
		userMessage:  message,
		kind:         "",
		errorMessage: message,
		timestamp:    time.Now(),
		status:       http.StatusInternalServerError,
		severity:     gLog.DefaultSeverity,
		stackTrace:   stacktrace.New(),
	}
	err.prev = previousError
	err.errorMessage = fmt.Errorf("%s: %w", message, previousError).Error()

	errI := Build(err, previousError)

	return errI
}

// Err creates a new Error.
func Err(message string) Error {
	return &FullError{
		userMessage:  message,
		kind:         "",
		errorMessage: message,
		timestamp:    time.Now(),
		status:       http.StatusInternalServerError,
		severity:     gLog.DefaultSeverity,
		stackTrace:   stacktrace.New(),
	}
}

// FullError implements the error interface.
// It will return the "developer" message in opposition to the user message,
// which is returned by FullError.Message
func (e *FullError) Error() string {
	return e.errorMessage
}

// Unwrap implements the errors.Unwrap interface
func (e *FullError) Unwrap() error {
	return e.prev
}

// StatusCode implements server.WithStatusCode
func (e *FullError) StatusCode() int {
	return e.status
}

// Message returns the user message.
func (e *FullError) Message() string {
	return e.userMessage
}

// Kind returns the error kind.
func (e *FullError) Kind() string {
	return e.kind
}

// Severity implements the gLog.WithSeverity interface.
func (e *FullError) Severity() gLog.Severity {
	if e.severity == gLog.DefaultSeverity {
		return gLog.ErrorSeverity
	}

	return e.severity
}

// Timestamp returns the error timestamp.
func (e *FullError) Timestamp() time.Time {
	return e.timestamp
}

// StackTrace returns the error stacktrace.
func (e *FullError) StackTrace() gLog.StackTrace {
	return e.stackTrace
}

// WithStatus sets the error status.
// It will also modify the severity for status >= 400
func (e *FullError) WithStatus(status int) Error {
	e.status = status

	if status >= 400 && status < 500 {
		e.severity = gLog.WarnSeverity
	} else if status >= 500 {
		e.severity = gLog.ErrorSeverity
	}

	return e
}

// WithMessage sets the user message.
func (e *FullError) WithMessage(message string) Error {
	e.userMessage = message

	return e
}

// WithKind sets the error kind.
func (e *FullError) WithKind(kind string) Error {
	e.kind = kind

	return e
}

// WithSeverity sets the error severity.
func (e *FullError) WithSeverity(severity gLog.Severity) Error {
	e.severity = severity

	return e
}

// WithTimestamp sets the error timestamp.
func (e *FullError) WithTimestamp(t time.Time) Error {
	e.timestamp = t

	return e
}

// WithStackTrace sets the error stacktrace.
func (e *FullError) WithStackTrace(trace gLog.StackTrace) Error {
	e.stackTrace = trace

	return e
}

// HttpError is used to json.Marshal or xml.Marshal FullError.
// You can use it to decode an incoming error.
type HttpError struct {
	Message string `json:"message" xml:"message"`
	Kind    string `json:"kind" xml:"kind"`
}

// MarshalJSON implements the json.Marshaler interface.
func (e *FullError) MarshalJSON() ([]byte, error) {
	return json.Marshal(HttpError{
		Message: e.userMessage,
		Kind:    e.kind,
	})
}

// MarshalXML implements the xml.Marshaler interface.
func (e *FullError) MarshalXML(encoder *xml.Encoder, start xml.StartElement) error {
	return encoder.EncodeElement(HttpError{
		Message: e.userMessage,
		Kind:    e.kind,
	}, start)
}
