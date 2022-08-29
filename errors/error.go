package errors

import (
	"fmt"
	"net/http"
)

// Error /
type Error struct {
	Msg    string `json:"message"`
	Code   int    `json:"code"`
	Origin error  `json:"-"`
}

// Error implements the error interface.
func (e Error) Error() string {
	return e.Message()
}

// Message returns the Error message. If the Msg field is not filled,
// try to call the origin Error method instead.
func (e Error) Message() string {
	if e.Msg == "" && e.Origin != nil {
		return e.Origin.Error()
	}
	return e.Msg
}

// StatusCode /
func (e Error) StatusCode() int {
	return e.Code
}

// Unwrap /
func (e Error) Unwrap() error {
	return e.Origin
}

type errOpt func(gapiError *Error)

// StatusCodeOpt override error status code /
func StatusCodeOpt(statusCode int) errOpt {
	return func(gapiError *Error) {
		gapiError.Code = statusCode
	}
}

// Wrap returns a new BadRequest Error. Use origin
// optional parameter to initialize the origin of this error.
func Wrap(origin error, message string, opts ...errOpt) Error {
	ge := Error{
		Code:   http.StatusInternalServerError,
		Msg:    message,
		Origin: origin,
	}

	for _, e := range opts {
		e(&ge)
	}

	return ge
}

// Errorf returns a new BadRequest Error with customer message.
func Errorf(statusCode int, format string, args ...interface{}) Error {
	return Error{
		Code: statusCode,
		Msg:  fmt.Sprintf(format, args...),
	}
}
