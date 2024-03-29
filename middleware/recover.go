package middleware

import (
	"net/http"

	"github.com/mwm-io/gapi/errors"
	"github.com/mwm-io/gapi/handler"
)

// Recover middleware will recover from panics and return the panic details as error.
type Recover struct{}

// Wrap implements the request.Middleware interface
func (r Recover) Wrap(h handler.Handler) handler.Handler {
	return handler.Func(func(w http.ResponseWriter, r *http.Request) (result interface{}, err error) {
		defer func() {
			if rec := recover(); rec != nil {
				err = errors.InternalServerError("panic", "Panic: %v", rec)
			}
		}()

		return h.Serve(w, r)
	})
}
