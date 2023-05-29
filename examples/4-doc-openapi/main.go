package main

import (
	"context"

	gLog "github.com/mwm-io/gapi/log"
	"github.com/mwm-io/gapi/server"

	"github.com/mwm-io/gapi/examples/4-doc-openapi/internal"
)

// This simple example show you how to make a CRUD for the users model.
// Each handler has a type.
func main() {
	ctx := context.Background()
	r := server.NewMux()

	// All following handler are complete example to know how to make a handler with request-scope data
	// (ie: body, path parameters, query parameters ...) and associated generated documentation based on
	// your code.

	// We chose this approach to secure an accurate documentation base and allow developers
	// to spend more time on enriching his API documentation.
	server.AddHandlerFactory(r, "GET", "/users", internal.SearchHandler)
	server.AddHandlerFactory(r, "POST", "/users", internal.CreateHandler)
	server.AddHandlerFactory(r, "GET", "/users/{id}", internal.GetOneHandler)
	server.AddHandlerFactory(r, "PUT", "/users/{id}", internal.UpdateHandler)
	server.AddHandlerFactory(r, "DELETE", "/users/{id}", internal.DeleteHandler)

	// server.AddDocHandlers add handler to expose API documentation.
	// Go to http://localhost:8080 to see the result
	if err := server.AddDocHandlers(r); err != nil {
		gLog.Error(ctx, err.Error())
	}

	gLog.Info(ctx, "Starting http server")

	if err := server.ServeAndHandleShutdown(r); err != nil {
		gLog.Emergency(ctx, err.Error())
	}

	gLog.Info(ctx, "Server stopped")
}
