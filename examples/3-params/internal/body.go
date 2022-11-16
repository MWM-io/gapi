package internal

import (
	"net/http"

	"github.com/mwm-io/gapi/handler"
	"github.com/mwm-io/gapi/middleware"
)

// NewBodyHandler decode the given body, store it in h.body and reply with the decoded boyd
func NewBodyHandler() handler.Handler {
	h := bodyHandler{}
	h.MiddlewareList = []handler.Middleware{
		middleware.JsonBody(&h.body),
		// You can use middleware.Body(&h.body) if you want to parse body with other content type
	}

	return &h
}

type bodyHandler struct {
	handler.WithMiddlewares

	body struct {
		First  string `json:"first"`
		Second string `json:"second"`
	}
}

// Serve implements handler.Handler and is the function called when a request is handled
func (h bodyHandler) Serve(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return h.body, nil
}
