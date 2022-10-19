package middleware

import (
	"fmt"
	"net/http"

	"github.com/mwm-io/gapi/errors"
	gLog "github.com/mwm-io/gapi/log"
	"github.com/mwm-io/gapi/server"
)

// Recover middleware will recover from panics and return the panic details as error.
type Recover struct{}

// Wrap implements the request.Middleware interface
func (r Recover) Wrap(h server.Handler) server.Handler {
	return server.HandlerFunc(func(w http.ResponseWriter, r *http.Request) (result interface{}, err error) {
		defer func() {
			if rec := recover(); rec != nil {
				err = errors.Wrap(fmt.Errorf("%v", rec), "Panic: %v").
					WithStatus(http.StatusInternalServerError).
					WithSeverity(gLog.CriticalSeverity)
			}
		}()

		return h.Serve(w, r)
	})
}
