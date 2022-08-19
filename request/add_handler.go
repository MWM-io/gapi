package request

import (
	"github.com/gorilla/mux"
)

// AddHandler add a new handler to mux router for given method and path
func AddHandler(router *mux.Router, method, path string, f HandlerFactory) {
	router.Methods(method).
		Path(path).
		Handler(httpHandler{factory: f})
}
