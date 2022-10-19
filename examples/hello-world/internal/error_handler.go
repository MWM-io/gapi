package internal

import (
	"net/http"

	"github.com/mwm-io/gapi/errors"
	"github.com/mwm-io/gapi/middleware"
	"github.com/mwm-io/gapi/server"
)

// ErrorHelloWorldHandler :
type ErrorHelloWorldHandler struct {
	server.MiddlewareHandler
}

// ErrorHelloWorldHandlerF /
func ErrorHelloWorldHandlerF() server.HandlerFactory {
	return func() server.Handler {
		return ErrorHelloWorldHandler{
			MiddlewareHandler: middleware.Core(),
		}
	}
}

// Serve /
func (h ErrorHelloWorldHandler) Serve(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	err := errors.Err("oups")

	return nil, err
}
