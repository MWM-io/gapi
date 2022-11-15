package openapi

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/swaggest/openapi-go/openapi3"

	"github.com/mwm-io/gapi/handler"
)

// PopulateReflector will add all the router routes into the given openapi3.Reflector
// You can add ignoredPath to ignore some registered routes.
func PopulateReflector(reflector *openapi3.Reflector, r *mux.Router, ignoredPaths []string) error {
	return r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, errPath := route.GetPathTemplate()
		if errPath != nil {
			return nil
		}

		for _, ignoredPath := range ignoredPaths {
			if path == ignoredPath {
				return nil
			}
		}

		methods, errMethods := route.GetMethods()
		if errMethods != nil {
			return nil
		}

		for _, method := range methods {
			if err := BuildOperation(reflector, route.GetHandler(), method, path); err != nil {
				return err
			}
		}

		return nil
	})
}

// BuildOperation adds the given h to the openapi3.Reflector
func BuildOperation(reflector *openapi3.Reflector, h interface{}, method, path string) error {
	docBuilder := NewDocBuilder(reflector, method, path)

	if handlerWithMiddlewares, ok := h.(handler.MiddlewareAware); ok {
		middlewareList := handlerWithMiddlewares.Middlewares()

		for _, middleware := range middlewareList {
			if middlewareWithDoc, isDocumented := middleware.(Documented); isDocumented {
				if err := middlewareWithDoc.Doc(docBuilder); err != nil {
					return err
				}
			}
		}
	}

	if handlerDoc, ok := h.(Documented); ok {
		if err := handlerDoc.Doc(docBuilder); err != nil {
			return err
		}
	}

	return docBuilder.
		WithError(http.StatusInternalServerError, "internal_error", "Internal server error, retry later or contact a developer if the problem persist").
		Commit().
		Error()
}
