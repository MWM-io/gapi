package middleware

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/mwm-io/gapi/handler"
	gLog "github.com/mwm-io/gapi/log"
)

// Log is a middleware that will:
// - set the given logger into the request's context.
// - log any error returned by the next handler
type Log struct {
	Logger *zap.Logger
}

// Wrap implements the request.Middleware interface
func (m Log) Wrap(h handler.Handler) handler.Handler {
	return handler.Func(func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
		if m.Logger == nil {
			m.Logger = gLog.Logger()
		}

		ctx := gLog.NewContext(r.Context(), m.Logger)

		resp, err := h.Serve(w, r.WithContext(ctx))

		if err != nil {
			m.Logger.Error(err.Error())
		}

		return resp, err
	})
}
