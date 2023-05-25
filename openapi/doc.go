/*
Package openapi provides helper functions to use github.com/swaggest/openapi-go/openapi3 more easily.

It also provides you with a set of rapidoc handlers to serve a live documentation of your API with rapidoc interface.

## How to add the rapidoc handlers

	import (
		"github.com/gorilla/mux"
		"github.com/mwm-io/gapi/openapi"
	)

	// your router with routes.
	var r *mux.Router

	err := openapi.AddRapidocHandlers(r, openapi.config{})
	if err != nil {
		log.Printf("error while adding rapidoc %+v\n", err)
	}

## How to document your handlers.

	// your handler
	type MyHandler struct{}

	// Doc implements the openapi.Documented interface
	func (m PathParameters) Doc(builder *openapi.DocBuilder) error {
		return builder.WithDescription("my handler description").
			WithParams({
				ObjectID string `path:"objectID"`
				IsFull   bool   `query:"is_full"`
			}).
			Error()
	}
*/
package openapi
