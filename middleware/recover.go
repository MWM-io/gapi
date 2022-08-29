package middleware

import (
	"net/http"

	"github.com/mwm-io/gapi/errors"
	"github.com/mwm-io/gapi/request"
)

// Recover middleware will recover from panics and return the panic details as error.
type Recover struct{}

// Wrap implements the request.Middleware interface
func (r Recover) Wrap(h request.Handler) request.Handler {
	return request.HandlerFunc(func(w http.ResponseWriter, r *http.Request) (result interface{}, err error) {
		defer func() {
			if rec := recover(); rec != nil {
				err = errors.Errorf(http.StatusInternalServerError, "Panic: %v", rec)
			}
		}()

		return h.Serve(w, r)
	})
}
