package middleware

import (
	"context"
	"net/http"

	"github.com/mwm-io/gapi/server"
)

// Logger is a logger able to log anything with a context.
type Logger interface {
	LogAnyC(context.Context, interface{})
}

// Log is a pre-processor that will set the request parameters into the Parameters field.
type Log struct {
	Logger Logger
}

// Wrap implements the request.Middleware interface
func (m Log) Wrap(h server.Handler) server.Handler {
	return server.HandlerFunc(func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
		resp, err := h.Serve(w, r)

		if err != nil {
			m.Logger.LogAnyC(r.Context(), err)
		}

		return resp, err
	})
}
