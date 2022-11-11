package hello

import (
	"net/http"

	"github.com/mwm-io/gapi/handler"
	"github.com/mwm-io/gapi/middleware"
)

// MakeJSONResponseHandler return a Hello world JSON encoded string
func MakeJSONResponseHandler() handler.Handler {
	return jsonResponseHandler{
		WithMiddlewares: handler.WithMiddlewares{
			MiddlewareList: []handler.Middleware{
				middleware.MakeResponseWriter().
					SetForcedContentType("application/json"),
			},
		},
	}
}

type jsonResponseHandler struct {
	handler.WithMiddlewares
}

// Serve implements handler.Handler and is the function called when a request is handled
func (h jsonResponseHandler) Serve(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return "Hello world", nil
}
