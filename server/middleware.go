package server

// Middleware is able to wrap a given Handler to build a new one.
type Middleware interface {
	Wrap(Handler) Handler
}

// MiddlewareAware is a struct containing its own middlewares.
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

// MiddlewareH returns a new MiddlewareHandler with the given middlewares.
func MiddlewareH(middlewares ...Middleware) MiddlewareHandler {
	return MiddlewareHandler{
		middlewares: middlewares,
	}
}

// HandlerFuncWithMiddleware is a handlerFunc with Middlewares.
// Use HandlerF to build a new one.
type HandlerFuncWithMiddleware struct {
	HandlerFunc
	MiddlewareHandler
}

// HandlerF builds a new HandlerFunc with middlewares.
func HandlerF(f HandlerFunc, middlewares ...Middleware) *HandlerFuncWithMiddleware {
	return &HandlerFuncWithMiddleware{
		HandlerFunc:       f,
		MiddlewareHandler: MiddlewareH(middlewares...),
	}
}
