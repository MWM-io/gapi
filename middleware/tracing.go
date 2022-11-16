package middleware

import (
	"net/http"

	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace/propagation"

	"github.com/mwm-io/gapi/handler"
)

// Tracing will add tracing information to the request context using opencensus..
type Tracing struct {
	// Propagation is the propagation format (how to read incoming trace information in the request).
	Propagation      propagation.HTTPFormat
	IsPublicEndpoint bool
}

// Wrap implements the request.Middleware interface
func (m Tracing) Wrap(h handler.Handler) handler.Handler {
	return handler.Func(func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
		var resp interface{}
		var err error

		newHandler := ochttp.Handler{
			Propagation: m.Propagation,
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				resp, err = h.Serve(w, r)

			}),
			IsPublicEndpoint: m.IsPublicEndpoint,
		}

		newHandler.ServeHTTP(w, r)

		return resp, err
	})
}
