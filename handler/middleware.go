package handler

// Middleware is able to wrap a given Handler to build a new one.
type Middleware interface {
	Wrap(Handler) Handler
}

// SortableMiddleware embeds Middleware type with a Weight method
// in order to get sorted.
type SortableMiddleware interface {
	Middleware
	Weight() int
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

// ByWeight sorts a list of Middleware by Weight
type ByWeight []Middleware

func (b ByWeight) Len() int {
	return len(b)
}

func (b ByWeight) Less(i, j int) bool {
	var iw, jw int

	if m, ok := b[i].(SortableMiddleware); ok {
		iw = m.Weight()
	}

	if m, ok := b[j].(SortableMiddleware); ok {
		jw = m.Weight()
	}

	return iw < jw
}

func (b ByWeight) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
