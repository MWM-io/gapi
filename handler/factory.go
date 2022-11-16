package handler

// Factory is a function that return a new Handler
// It is useful if you want to create a Handler that will carry request-scoped data.
type Factory func() Handler

// // ServeHTTP implements the http.Handler interface.
// // It will try to find the first Handler in the chain that implements the MiddlewareAware interface to add the middlewares.
// func (h Factory) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	handlerInstance := h()
//
// 	middlewares := h.Middlewares(handlerInstance)
// 	for i := len(middlewares) - 1; i >= 0; i-- {
// 		handlerInstance = middlewares[i].Wrap(handlerInstance)
// 	}
//
// 	_, _ = handlerInstance.Serve(w, r)
// }

// // Middlewares implements the MiddlewareAware interface.
// func (h Factory) Middlewares(handlerInstance Handler) []Middleware {
// 	middlewareList := make([]Middleware, len(server.DefaultMiddlewares))
// 	copy(middlewareList, server.DefaultMiddlewares)
//
// 	middlewareHandler, ok := handlerInstance.(MiddlewareAware)
// 	if !ok {
// 		return nil
// 	}
//
// 	return append(middlewareList, middlewareHandler.Middlewares()...)
// }
