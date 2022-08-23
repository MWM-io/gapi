package internal

import (
	"github.com/mwm-io/gapi/error"
	"github.com/mwm-io/gapi/request"
)

// JsonHelloWorldHandler :
type JsonHelloWorldHandler struct{}

// JsonHelloWorldHandlerF /
func JsonHelloWorldHandlerF() request.HandlerFactory {
	return func() request.Handler {
		return JsonHelloWorldHandler{}
	}
}

// Serve /
func (h JsonHelloWorldHandler) Serve(_ request.WrappedRequest) (interface{}, error.Error) {
	return struct{ Text string }{Text: "Hello World"}, nil
}
