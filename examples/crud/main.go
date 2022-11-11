package main

import (
	"log"

	gLog "github.com/mwm-io/gapi/log"
	"github.com/mwm-io/gapi/server"

	"github.com/mwm-io/gapi/examples/crud/internal"
)

// This simple example show you how to make a CRUD for the users model.
// Each handler has a type.
func main() {
	r := server.NewMux()

	// See the ListHandler for more information on how to make a handler with request-scope data
	// (ie: body, path parameters, query parameters ...)
	server.AddHandlerFactory(r, "GET", "/users", internal.ListHandlerF)
	server.AddHandlerFactory(r, "POST", "/users", internal.PostHandlerF)
	server.AddHandlerFactory(r, "GET", "/users/{id}", internal.GetHandlerF)
	server.AddHandlerFactory(r, "PUT", "/users/{id}", internal.PutHandlerF)
	server.AddHandlerFactory(r, "DELETE", "/users/{id}", internal.DeleteHandlerF)

	if err := server.AddDocHandlers(r); err != nil {
		log.Printf("error while adding rapidoc %+v\n", err)
	}

	gLog.Info("Starting http server")

	if err := server.ServeAndHandleShutdown(r); err != nil {
		gLog.LogAny(err)
	}

	gLog.Info("Server stopped")
}
