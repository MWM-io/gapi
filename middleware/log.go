package middleware

import (
	"context"
	"net/http"

	"github.com/mwm-io/gapi/handler"
	gLog "github.com/mwm-io/gapi/log"
)

// Log is a middleware that will:
// - populate the logger context with the request logger
// - set the given logger into the request's context.
// - log any error returned by the next handler
type Log struct {
	Logger *gLog.Logger
}

// Wrap implements the request.Middleware interface
func (m Log) Wrap(h handler.Handler) handler.Handler {
	return handler.Func(func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
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
