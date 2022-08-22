package internal

import (
	"github.com/mwm-io/gapi/request"
	"github.com/mwm-io/gapi/response"
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
func (h JsonHelloWorldHandler) Serve(_ request.WrappedRequest) (interface{}, response.Error) {
	return struct{ Text string }{Text: "Hello World"}, nil
}
