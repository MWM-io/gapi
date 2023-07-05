package err

import (
	"net/http"

	"github.com/mwm-io/gapi/errors"
	"github.com/mwm-io/gapi/handler"
	"github.com/mwm-io/gapi/middleware"
)

// MakeAutoResponseHandler return an encoded error according to Accept header
func MakeAutoResponseHandler() handler.Handler {
	return autoResponseHandler{
		WithMiddlewares: handler.WithMiddlewares{
			MiddlewareList: []handler.Middleware{
				middleware.MakeResponseWriter(),
			},
		},
	}
}

type autoResponseHandler struct {
	handler.WithMiddlewares
}

// Serve implements handler.Handler and is the function called when a request is handled
func (h autoResponseHandler) Serve(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return nil, errors.Err("example", "Hello World").
		WithStatus(http.StatusTeapot)
}
