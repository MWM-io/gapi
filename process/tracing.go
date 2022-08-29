package process

import (
	"net/http"

	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace/propagation"
	
	"github.com/mwm-io/gapi/request"
)

// Tracing will add tracing information
type Tracing struct {
	Propagation      propagation.HTTPFormat
	IsPublicEndpoint bool
}

// Wrap implements the request.Middleware interface
func (m Tracing) Wrap(h request.Handler) request.Handler {
	return request.HandlerFunc(func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
		var resp interface{}
		var err error

		h := ochttp.Handler{
			Propagation: m.Propagation,
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				resp, err = h.Serve(w, r)

			}),
			IsPublicEndpoint: m.IsPublicEndpoint,
		}

		h.ServeHTTP(w, r)

		return resp, err
	})
}
