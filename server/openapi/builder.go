package openapi

import (
	"github.com/gorilla/mux"
	"github.com/swaggest/openapi-go/openapi3"

	"github.com/mwm-io/gapi/server"
)

// PopulateReflector will add all the router routes into the given openapi3.Reflector
// You can add ignoredPath to ignore some of the registered routes.
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

// BuildOperation adds the given handler to the openapi3.Reflector
func BuildOperation(reflector *openapi3.Reflector, handler interface{}, method, path string) error {
	if httpHandler, ok := handler.(server.HttpHandler); ok {
		return BuildOperation(reflector, httpHandler.Handler, method, path)
	}

	if handlerFactory, ok := handler.(server.HandlerFactory); ok {
		return BuildOperation(reflector, handlerFactory(), method, path)
	}

	docBuilder := NewOperationBuilder(reflector, method, path)

	if handlerWithMiddlewares, ok := handler.(server.MiddlewareAware); ok {
		for _, middleware := range handlerWithMiddlewares.Middlewares() {
			if middlewareWithDoc, ok := middleware.(OperationDescriptor); ok {
				if err := middlewareWithDoc.Doc(docBuilder); err != nil {
					return err
				}
			}
		}
	}

	if handlerDoc, ok := handler.(OperationDescriptor); ok {
		if err := handlerDoc.Doc(docBuilder); err != nil {
			return err
		}
	}

	return docBuilder.
		Commit().
		Error()
}
