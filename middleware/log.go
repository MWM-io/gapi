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
type Log struct{}

// Wrap implements the request.Middleware interface
func (m Log) Wrap(h handler.Handler) handler.Handler {
	return handler.Func(func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
		ctx := r.Context()
		l := gLog.Logger(ctx)

		ctx = gLog.NewRefContext(ctx, l)
		r = r.WithContext(gLog.NewContext(ctx, l))

		resp, err := h.Serve(w, r)

		if err != nil {
			latest := gLog.LatestLogger(r.Context())
			errLog := &gLog.Log{}
			errLog.SetFunc(latest.WithOptions(zap.AddCallerSkip(1)).Error)
			errLog.LogError(err)
		}

		return resp, err
	})
}
