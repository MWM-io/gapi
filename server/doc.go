/*
Package server provides a simple router based https://github.com/gorilla/mux

	r := router.NewMux()

	// add your routes with the request package

	err := ServeAndHandleShutdown(r)
	if err != nil {
		log.Fatal(err)
	}

*/
package server
