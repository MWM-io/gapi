package request

import (
	"net/http"

	"github.com/gorilla/mux"
)

// AddHandler add a new handler factory to mux router for given method and path
func AddHandler(router *mux.Router, method, path string, f HandlerFactory) {
	router.Methods(method).
		Path(path).
		Handler(f)
}

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

// HandlerFactory is a function that return a new Handler
// It is useful if you want to create a Handler that will carry request-scoped data.
type HandlerFactory func() Handler

// Serve implements the Handler interface
func (h HandlerFactory) Serve(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return h().Serve(w, r)
}

func (h HandlerFactory) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler := h()

	if middlewareHandler, ok := handler.(MiddlewareAware); ok {
		middlewares := middlewareHandler.Middlewares()
		for i := len(middlewares) - 1; i >= 0; i-- {
			handler = middlewares[i].Wrap(handler)
		}
	}

	_, _ = handler.Serve(w, r)
}

// Middlewares implements the MiddlewareAware interface.
func (h HandlerFactory) Middlewares() []Middleware {
	handler := h()

	if middlewareHandler, ok := handler.(MiddlewareAware); ok {
		return middlewareHandler.Middlewares()
	}

	return nil
}
