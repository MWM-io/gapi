package internal

import (
	"net/http"

	"github.com/mwm-io/gapi/middleware"
	"github.com/mwm-io/gapi/server"
)

// JsonHelloWorldHandler :
type JsonHelloWorldHandler struct {
	server.MiddlewareHandler
}

// JsonHelloWorldHandlerF /
func JsonHelloWorldHandlerF() server.HandlerFactory {
	return func() server.Handler {
		return JsonHelloWorldHandler{
			MiddlewareHandler: middleware.Core(),
		}
	}
}

// Serve /
func (h JsonHelloWorldHandler) Serve(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return struct{ Text string }{Text: "Hello World"}, nil
}
