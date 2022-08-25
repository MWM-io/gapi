package internal

import (
	"net/http"

	"github.com/mwm-io/gapi/process"
	"github.com/mwm-io/gapi/request"
)

// JsonHelloWorldHandler :
type JsonHelloWorldHandler struct {
	request.MiddlewareHandler
}

// JsonHelloWorldHandlerF /
func JsonHelloWorldHandlerF() request.HandlerFactory {
	return func() request.Handler {
		return JsonHelloWorldHandler{
			MiddlewareHandler: process.Core(),
		}
	}
}

// Serve /
func (h JsonHelloWorldHandler) Serve(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return struct{ Text string }{Text: "Hello World"}, nil
}
