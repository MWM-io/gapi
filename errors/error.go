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

	WithMessage(format string, args ...interface{}) Error
	WithKind(string) Error
	WithStatus(int) Error
}

// FullError is a concrete error that implements the Error interface
type FullError struct {
	userMessage  string
	kind         string
	errorMessage string
	status       int
	timestamp    time.Time
	sourceErr    error
}

// Wrap will wrap the given error and return a new Error.
func Wrap(err error) Error {
	if err == nil {
		return nil
	}

	newErr := &FullError{
		userMessage:  err.Error(),
		kind:         "",
		errorMessage: err.Error(),
		timestamp:    time.Now(),
		status:       http.StatusInternalServerError,
		sourceErr:    err,
	}

	return newErr
}

// Err creates a new Error.
func Err(format string, args ...interface{}) Error {
	message := fmt.Sprintf(format, args...)

	return &FullError{
		userMessage:  message,
		kind:         "",
		errorMessage: message,
		timestamp:    time.Now(),
		status:       http.StatusInternalServerError,
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
