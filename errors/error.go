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

// ErrorI represents the interface for the Error struct.
// It is necessary so we can compare nil errors. (nil.(*Error) != nil)
type ErrorI interface {
	error
	Unwrap() error
	json.Marshaler
	xml.Marshaler
	StatusCode() int
	Severity() gLog.Severity
	StackTrace() gLog.StackTrace
	WithStatus(status int) ErrorI
	WithMessage(message string) ErrorI
	WithKind(kind string) ErrorI
	WithSeverity(severity gLog.Severity) ErrorI
}

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

func Wrap(previousError error, message string) ErrorI {
	if previousError == nil {
		return nil
	}

	err := Error{
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

	errI := Build(err.ErrorI(), previousError)

	if errI.Severity() == gLog.DefaultSeverity {
		if status := errI.StatusCode(); status >= 400 && status < 500 {
			errI = errI.WithSeverity(gLog.WarnSeverity)
		} else {
			errI = errI.WithSeverity(gLog.ErrorSeverity)
		}
	}

	return errI
}

func Err(message string) ErrorI {
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

func (e Error) ErrorI() ErrorI {
	return e
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

func (e Error) WithStatus(status int) ErrorI {
	e.status = status

	return e
}

func (e Error) WithMessage(message string) ErrorI {
	e.userMessage = message

	return e
}

func (e Error) WithKind(kind string) ErrorI {
	e.kind = kind

	return e
}

func (e Error) WithSeverity(severity gLog.Severity) ErrorI {
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
