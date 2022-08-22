package internal

import (
	"errors"

	"github.com/mwm-io/gapi/request"
	"github.com/mwm-io/gapi/response"
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
func (h ErrorHelloWorldHandler) Serve(_ request.WrappedRequest) (interface{}, response.Error) {
	err := response.GapiError{
		Msg:    "Hello World",
		Code:   500,
		Origin: errors.New("something went wrong"),
	}

	return nil, err
}
