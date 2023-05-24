package server

import (
	"sync"

	"github.com/mwm-io/gapi/handler"
	"github.com/mwm-io/gapi/middleware"
)

var middlewareMU sync.Mutex

// UseMiddlewares appends a given list of handler.Middleware to middlewares chain.
//
// Middleware can be used to intercept or otherwise modify requests and/or responses, and
// are executed in list order.
func UseMiddlewares(middlewares ...handler.Middleware) {
	middlewareMU.Lock()
	defer middlewareMU.Unlock()

	middleware.Defaults = append(middleware.Defaults, middlewares...)
}
