package openapi

import (
	"encoding/json"
	"net/http"

	"github.com/swaggest/openapi-go/openapi3"
)

// Handler is a server.Handler that will return a json openapi definition of the given reflector.
type Handler struct {
	reflector *openapi3.Reflector
	auth      Authorization
}

func NewHandler(reflector *openapi3.Reflector, auth Authorization) *Handler {
	return &Handler{
		reflector: reflector,
		auth:      auth,
	}
}

// Serve implements the server.Handler interface
func (h *Handler) Serve(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	if h.auth != nil {
		authorized, err := h.auth.Authorize(w, r)
		if err != nil {
			return nil, err
		}
		if !authorized {
			return h.auth.Login(w, r)
		}
	}

	body, err := json.Marshal(h.reflector.Spec)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	_, err = w.Write(body)
	return nil, err
}
