package handler

// Middleware is able to wrap a given Handler to build a new one.
type Middleware interface {
	Wrap(Handler) Handler
}

// MiddlewareAware is a struct containing its own middlewares.
type MiddlewareAware interface {
	Middlewares() []Middleware
}

// WithMiddlewares implements the request.MiddlewareAware interface
// and is meant to be embedded into the final handler
type WithMiddlewares struct {
	MiddlewareList []Middleware
}

// Middlewares implements the request.MiddlewareAware interface
func (p WithMiddlewares) Middlewares() []Middleware {
	return p.MiddlewareList
}
