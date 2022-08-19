package internal

import (
	"github.com/mwm-io/gapi/request"
	"github.com/mwm-io/gapi/response"
)

// HelloWorldHandler :
type HelloWorldHandler struct {
}

// HelloWorldHandlerF /
func HelloWorldHandlerF() request.HandlerFactory {
	return func() request.Handler {
		return HelloWorldHandler{}
	}
}

// Serve /
func (h HelloWorldHandler) Serve(w request.WrappedRequest) (interface{}, response.Error) {
	return "Hello World", nil
}
