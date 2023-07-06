package errors

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"time"
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
	Timestamp() time.Time
	CallerName() string
	Caller() string
	Callstack() []string

	WithMessage(format string, args ...interface{}) Error
	WithKind(string) Error
	WithStatus(int) Error
	WithError(error) Error
}

// FullError is a concrete error that implements the Error interface
type FullError struct {
	userMessage  string
	kind         string
	errorMessage string
	status       int
	timestamp    time.Time
	sourceErr    error
	callerName   string
	caller       string
	callstack    []string
}

// Wrap will wrap the given error and return a new Error.
func Wrap(err error) Error {
	if err == nil {
		return nil
	}

	if castedErr, ok := err.(Error); ok {
		return castedErr
	}

	for _, builder := range errorBuilders {
		if gErr := builder(err); gErr != nil {
			return gErr
		}
	}

	callerName, caller, callstack := GetCallers()

	newErr := &FullError{
		userMessage:  err.Error(),
		kind:         "internal_error",
		errorMessage: err.Error(),
		status:       http.StatusInternalServerError,
		timestamp:    time.Now(),
		sourceErr:    err,
		callerName:   callerName,
		caller:       caller,
		callstack:    callstack,
	}

	return newErr
}

// Err creates a new Error.
func Err(kind, format string, args ...interface{}) Error {
	message := fmt.Sprintf(format, args...)

	callerName, caller, callstack := GetCallers()

	return &FullError{
		userMessage:  message,
		kind:         kind,
		errorMessage: message,
		timestamp:    time.Now(),
		status:       http.StatusInternalServerError,
		callerName:   callerName,
		caller:       caller,
		callstack:    callstack,
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
	return e.sourceErr
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

// Timestamp returns the error timestamp.
func (e *FullError) Timestamp() time.Time {
	return e.timestamp
}

// WithStatus sets the error status.
// It will also modify the severity for status >= 400
func (e *FullError) WithStatus(status int) Error {
	e.status = status

	return e
}

// WithMessage sets the user message.
func (e *FullError) WithMessage(format string, args ...interface{}) Error {
	e.userMessage = fmt.Sprintf(format, args...)

	return e
}

// WithKind sets the error kind.
func (e *FullError) WithKind(kind string) Error {
	e.kind = kind

	return e
}

// WithError wrap source error.
func (e *FullError) WithError(err error) Error {
	e.sourceErr = err
	if e.errorMessage == "" {
		e.errorMessage = err.Error()
	}

	return e
}

// CallerName implements the error interface.
// It will return the name of the function that created the error
func (e *FullError) CallerName() string {
	return e.callerName
}

// Caller implements the error interface.
// It will return the name of the function that created the error
func (e *FullError) Caller() string {
	return e.callerName
}

// Callstack implements the error interface.
// It will return the complete callstack of the error creation
func (e *FullError) Callstack() []string {
	return e.callstack
}

// CallerName() string
// Caller() string
// Callstack() []string

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
