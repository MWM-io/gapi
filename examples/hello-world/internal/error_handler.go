package internal

import (
	"errors"

	"github.com/mwm-io/gapi/error"
	"github.com/mwm-io/gapi/request"
)

// ErrorHelloWorldHandler :
type ErrorHelloWorldHandler struct{}

// ErrorHelloWorldHandlerF /
func ErrorHelloWorldHandlerF() request.HandlerFactory {
	return func() request.Handler {
		return ErrorHelloWorldHandler{}
	}
}

// Serve /
func (h ErrorHelloWorldHandler) Serve(_ request.WrappedRequest) (interface{}, error.Error) {
	err := error.GapiError{
		Msg:    "Hello World",
		Code:   500,
		Origin: errors.New("something went wrong"),
	}

	return nil, err
}
