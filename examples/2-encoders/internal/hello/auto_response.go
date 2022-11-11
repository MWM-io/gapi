package hello

import (
	"net/http"

	"github.com/mwm-io/gapi/handler"
	"github.com/mwm-io/gapi/middleware"
)

// MakeAutoResponseHandler return a Hello world string encoded according to Accept header
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
	return "Hello world", nil
}
