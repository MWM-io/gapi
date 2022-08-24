package request

type Middleware interface {
	Wrap(Handler) Handler
}

type MiddlewareAware interface {
	Middlewares() []Middleware
}

// MiddlewareHandler implements the request.MiddlewareAware interface
// and is meant to be embedded into the final handler
type MiddlewareHandler struct {
	middlewares []Middleware
}

// Middlewares implements the request.MiddlewareAware interface
func (p MiddlewareHandler) Middlewares() []Middleware {
	return p.middlewares
}

// MiddlewareH returns a new MiddlewareHandler with the given preProcesses.
func MiddlewareH(middlewares ...Middleware) MiddlewareHandler {
	return MiddlewareHandler{
		middlewares: middlewares,
	}
}
