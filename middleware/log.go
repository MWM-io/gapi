package middleware

import (
	"context"
	"net/http"

	gLog "github.com/mwm-io/gapi/log"
	"github.com/mwm-io/gapi/server"
)

// Log is a pre-processor that will set the request parameters into the Parameters field.
type Log struct {
	Logger *gLog.Logger
}

// Wrap implements the request.Middleware interface
func (m Log) Wrap(h server.Handler) server.Handler {
	return server.HandlerFunc(func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
		var ctx context.Context
		if m.Logger == nil {
			ctx = gLog.NewContext(r.Context(), gLog.GlobalLogger())
		} else {
			ctx = gLog.NewContext(r.Context(), m.Logger)
		}

		ctx = gLog.CtxWithContext(ctx)
		r = r.WithContext(ctx)

		resp, err := h.Serve(w, r)

		if err != nil {
			gLog.LogAnyC(r.Context(), err)
		}

		return resp, err
	})
}
