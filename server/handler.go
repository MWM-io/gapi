package server

import (
	"net/http"
)

// Handler is able to respond to a http request
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
