package handler

import (
	"net/http"
)

// Func type is an adapter to allow the use of
// ordinary functions as HTTP handlers. If f is a function
// with the appropriate signature, handler.Func(f) is a
// Handler that calls f
type Func func(http.ResponseWriter, *http.Request) (interface{}, error)

// Serve implements the handler.Handler interface
func (f Func) Serve(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return f(w, r)
}
