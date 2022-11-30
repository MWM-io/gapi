package openapi

import (
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/swaggest/openapi-go/openapi3"

	"github.com/mwm-io/gapi/handler"
)

// SpecOpenAPIHandler is a server.SpecOpenAPIHandler that will return a json openapi definition of the given reflector.
type SpecOpenAPIHandler struct {
	handler.WithMiddlewares
	_reflector      *openapi3.Reflector
	reflectorOnce   sync.Once
	computeDocError error
	router          *mux.Router
}

// NewSpecOpenAPIHandler builds a new SpecOpenAPIHandler, serving the api definition from the openapi3.Reflector,
// and checking auth access with the given Authorization.
func NewSpecOpenAPIHandler(router *mux.Router, middlewares ...handler.Middleware) *SpecOpenAPIHandler {
	return &SpecOpenAPIHandler{
		router: router,
		WithMiddlewares: handler.WithMiddlewares{
			MiddlewareList: middlewares,
		},
	}
}

// getReflector return the reflector fill with generated documentation.
// With this way of doing the documentation is generated only at the first consultation.
// If we generate the documentation when creating the Handler the server startup would be impacted
func (h *SpecOpenAPIHandler) getReflector() (*openapi3.Reflector, error) {
	h.reflectorOnce.Do(func() {
		h._reflector = new(openapi3.Reflector)
		h._reflector.SpecEns().Info.WithTitle(Config.GetDocPageTitle())
		h.computeDocError = PopulateReflector(h._reflector, h.router, Config.ignoredPaths())
	})

	return h._reflector, h.computeDocError
}

// Serve implements the handler.Handler interface
func (h *SpecOpenAPIHandler) Serve(_ http.ResponseWriter, _ *http.Request) (interface{}, error) {
	reflector, err := h.getReflector()
	if err != nil {
		return nil, err
	}

	return reflector.Spec, nil
}
