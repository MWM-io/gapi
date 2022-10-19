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

type Error struct {
	userMessage  string
	kind         string
	errorMessage string
	status       int
	severity     gLog.Severity
	timestamp    time.Time
	stackTrace   stacktrace.StackTrace
	prev         error
}

func Wrap(previousError error, message string) *Error {
	if previousError == nil {
		return nil
	}

	err := Err(message)

	err = Build(err, previousError)
	err.prev = previousError
	err.errorMessage = fmt.Errorf("%s: %w", message, previousError).Error()

	return &err
}

func Err(message string) Error {
	return Error{
		userMessage:  message,
		kind:         "",
		errorMessage: message,
		timestamp:    time.Now(),
		status:       http.StatusInternalServerError,
		severity:     gLog.ErrorSeverity,
		stackTrace:   stacktrace.New(),
	}
}

// Error implements the error interface.
func (e Error) Error() string {
	return e.errorMessage
}

// Unwrap /
func (e Error) Unwrap() error {
	return e.prev
}

// StatusCode implements server.WithStatusCode
func (e Error) StatusCode() int {
	return e.status
}

func (e Error) Severity() gLog.Severity {
	return e.severity
}

func (e Error) StackTrace() gLog.StackTrace {
	return e.stackTrace
}

func (e Error) WithStatus(status int) Error {
	e.status = status

	return e
}

func (e Error) WithMessage(message string) Error {
	e.userMessage = message

	return e
}

func (e Error) WithKind(kind string) Error {
	e.kind = kind

	return e
}

func (e Error) WithSeverity(severity gLog.Severity) Error {
	e.severity = severity

	return e
}

type HttpError struct {
	Message string `json:"message" xml:"message"`
	Kind    string `json:"kind" xml:"kind"`
}

func (e Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(HttpError{
		Message: e.userMessage,
		Kind:    e.kind,
	})
}

func (e Error) MarshalXML(encoder *xml.Encoder, start xml.StartElement) error {
	return encoder.EncodeElement(HttpError{
		Message: e.userMessage,
		Kind:    e.kind,
	}, start)
}
