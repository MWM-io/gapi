package openapi

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/swaggest/openapi-go/openapi3"
)

// SpecOpenAPIHandler is a server.SpecOpenAPIHandler that will return a json openapi definition of the given reflector.
type SpecOpenAPIHandler struct {
	_reflector      *openapi3.Reflector
	reflectorOnce   sync.Once
	computeDocError error
	router          *mux.Router
	// auth      Authorization
}

// NewSpecOpenAPIHandler builds a new SpecOpenAPIHandler, serving the api definition from the openapi3.Reflector,
// and checking auth access with the given Authorization.
func NewSpecOpenAPIHandler(router *mux.Router) *SpecOpenAPIHandler {
	return &SpecOpenAPIHandler{
		router: router,
		// auth:      auth, TODO : change it by a middleware
	}
}

// getReflector return the reflector fill with generated documentation.
// With this way of doing the documentation is generated only at the first consultation.
// If we generate the documentation when creating the Handler the server startup would be impacted
func (h *SpecOpenAPIHandler) getReflector() (*openapi3.Reflector, error) {
	h.reflectorOnce.Do(func() {
		h._reflector = new(openapi3.Reflector)
		h.computeDocError = PopulateReflector(h._reflector, h.router, Config.ignoredPaths())
	})

	return h._reflector, h.computeDocError
}

// Serve implements the handler.Handler interface
func (h *SpecOpenAPIHandler) Serve(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	// if h.auth != nil {
	// 	authorized, err := h.auth.Authorize(w, r)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	if !authorized {
	// 		return h.auth.Login(w, r)
	// 	}
	// }

	reflector, err := h.getReflector()
	if err != nil {
		// TODO : change it by response writer middleware
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	body, errMarshal := json.Marshal(reflector.Spec)
	if errMarshal != nil {
		// TODO : change it by response writer middleware
		http.Error(w, errMarshal.Error(), http.StatusInternalServerError)
		return nil, errMarshal
	}

	// TODO : change it by response writer middleware
	_, err = w.Write(body)
	return nil, errMarshal
}
