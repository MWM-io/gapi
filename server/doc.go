/*
Package server provides a simple router based https://github.com/gorilla/mux,
with a Handler type return the response (not serialized) and an error,
as well as a Middleware type to allow middleware for this new Handler type.

You can see middleware implementation in the github.com/mwm-io/gapi/middleware package.


```go
	import (
		"http"
		"log"

		"github.com/mwm-io/gapi/server"
	)

	r := server.NewMux()

	// Add your http handlers.
	var h Handler
	server.AddHandler(r, http.MethodGet, "/hello", h)

	// Add your server options here.
	err := server.ServeAndHandleShutdown(r)
	if err != nil {
		log.Fatal(err)
	}
```

*/
package server
