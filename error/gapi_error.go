package error

import (
	"fmt"
	"net/http"
)

// GapiError /
type GapiError struct {
	Msg    string `json:"message"`
	Code   int    `json:"code"`
	Origin error  `json:"-"`
}

// Error implements the error interface.
func (e GapiError) Error() string {
	return e.Message()
}

// Message returns the Error message. If the Msg field is not filled,
// try to call the origin Error method instead.
func (e GapiError) Message() string {
	if e.Msg == "" && e.Origin != nil {
		return e.Origin.Error()
	}
	return e.Msg
}

// StatusCode /
func (e GapiError) StatusCode() int {
	return e.Code
}

// Unwrap /
func (e GapiError) Unwrap() error {
	return e.Origin
}

type errOpt func(gapiError *GapiError)

// StatusCodeOpt override error status code /
func StatusCodeOpt(statusCode int) errOpt {
	return func(gapiError *GapiError) {
		gapiError.Code = statusCode
	}
}

// Wrap returns a new BadRequest Error. Use origin
// optional parameter to initialize the origin of this error.
func Wrap(origin error, message string, opts ...errOpt) GapiError {
	ge := GapiError{
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
func Errorf(statusCode int, format string, args ...interface{}) GapiError {
	return GapiError{
		Code: statusCode,
		Msg:  fmt.Sprintf(format, args...),
	}
}
