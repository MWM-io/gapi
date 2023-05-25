package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mwm-io/gapi/handler"
	"github.com/mwm-io/gapi/middleware"
	"github.com/mwm-io/gapi/openapi"
)

// AddHandlerFactory register a new handler factory to the given mux router on a given method and path.
// You can use it if you want a new instance of your handler for each call.
// It must be use useful if :
//
//   - you use a middleware for handle request params like middleware.BodyDecoder, middleware.PathParameters, etc.
//   - you store properties in your handler struct during Serve process
func AddHandlerFactory(router *mux.Router, method, path string, f handler.Factory) {
	router.Methods(method).
		Path(path).
		Handler(defaultHandleEngine{
			getHandler: f,
		})
}

// AddHandler register a new handler to the given mux router on a given method and path.
func AddHandler(router *mux.Router, method, path string, f handler.Handler) {
	router.Methods(method).
		Path(path).
		Handler(defaultHandleEngine{
			getHandler: func() handler.Handler {
				return f
			},
		})
}

type defaultHandleEngine struct {
	getHandler func() handler.Handler
}

func (e defaultHandleEngine) getMiddlewareList(h handler.Handler) []handler.Middleware {
	middlewareList := make([]handler.Middleware, len(middleware.Defaults))
	copy(middlewareList, middleware.Defaults)

	// if current handler have custom middlewares, add them to the list
	if middlewareHandler, ok := h.(handler.MiddlewareAware); ok {
		middlewareList = append(middlewareList, middlewareHandler.Middlewares()...)
	}

	return middlewareList
}

// ServeHTTP is the function called by mux when a request is handled
func (e defaultHandleEngine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get handler to serve
	h := e.getHandler()

	// get middleware list for this handler
	middlewareList := e.getMiddlewareList(h)

	// Execute all middleware with list order
	for i := len(middlewareList) - 1; i >= 0; i-- {
		h = middlewareList[i].Wrap(h)
	}

	// Ignore response & error : must be handle by response writer middleware
	_, _ = h.Serve(w, r)
}

// Doc is the function called by openapi during doc generation
func (e defaultHandleEngine) Doc(builder *openapi.DocBuilder) error {
	// get handler to serve
	h := e.getHandler()

	// get middleware list for this handler
	middlewareList := e.getMiddlewareList(h)

	// Execute all middleware Doc func (if exist) with list order
	for i := len(middlewareList) - 1; i >= 0; i-- {
		documentedMiddleware, ok := middlewareList[i].(openapi.Documented)
		if !ok {
			continue
		}

		if err := documentedMiddleware.Doc(builder); err != nil {
			return err
		}
	}

	// Execute handler Doc func (if exist)
	if documentedHandler, ok := h.(openapi.Documented); ok {
		return documentedHandler.Doc(builder)
	}

	return nil
}
