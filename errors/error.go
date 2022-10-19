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

func Err(message string, previousError ...error) Error {
	err := Error{
		userMessage:  message,
		kind:         "",
		errorMessage: message,
		timestamp:    time.Now(),
		status:       http.StatusInternalServerError,
		severity:     gLog.ErrorSeverity,
		stackTrace:   stacktrace.New(),
	}

	if len(previousError) == 0 {
		return err
	}

	if len(previousError) > 1 {
		gLog.Critical("you cannot call errors.E with more than one error")
		return Error{}
	}

	err = Build(err, previousError[0])
	err.prev = previousError[0]
	err.errorMessage = fmt.Errorf("%s\n  %w", message, previousError[0]).Error()

	return err
}

func Warn(message string, previousError ...error) Error {
	return Err(message, previousError...).
		WithSeverity(gLog.WarnSeverity)
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
