package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Handler is able to respond to a http request and return the response it wants to write and an error.
// You need to use a middleware to write your response if you just return your response and not write it.
// (see github.com/mwm-io/gapi/middlewares.BodyUnmarshaler)
type Handler interface {
	Serve(http.ResponseWriter, *http.Request) (interface{}, error)
}

// HandlerFunc type is an adapter to allow the use of
// ordinary functions as HTTP handlers. If f is a function
// with the appropriate signature, HandlerFunc(f) is a
// Handler that calls f.
type HandlerFunc func(http.ResponseWriter, *http.Request) (interface{}, error)

// Serve implements the Handler interface.
func (h HandlerFunc) Serve(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return h(w, r)
}

// AddHandler add a new handler to the given mux router on a given method and path.
func AddHandler(router *mux.Router, method, path string, f Handler) {
	router.Methods(method).
		Path(path).
		Handler(HttpHandler{f})
}

// HttpHandler is a wrapper of Handler that implements the http.Handler interface.
type HttpHandler struct {
	Handler
}

// ServeHTTP implements the http.Handler interface.
// It will try to find the first Handler in the chain that implements the MiddlewareAware interface to add the middlewares.
func (h HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler := h.Handler

	middlewares := h.Middlewares()
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i].Wrap(handler)
	}

	_, _ = handler.Serve(w, r)
}

// Middlewares implements the MiddlewareAware interface.
func (h HttpHandler) Middlewares() []Middleware {
	middlewareHandler, ok := h.Handler.(MiddlewareAware)
	if !ok {
		return nil
	}

	return middlewareHandler.Middlewares()
}

// HandlerFactory is a function that return a new Handler
// It is useful if you want to create a Handler that will carry request-scoped data.
type HandlerFactory func() Handler

// Serve implements the Handler interface
func (h HandlerFactory) Serve(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	handler := h()

	if middlewareHandler, ok := handler.(MiddlewareAware); ok {
		middlewares := middlewareHandler.Middlewares()
		for i := len(middlewares) - 1; i >= 0; i-- {
			handler = middlewares[i].Wrap(handler)
		}
	}

	return handler.Serve(w, r)
}
