package main

import (
	gLog "github.com/mwm-io/gapi/log"
	"github.com/mwm-io/gapi/server"

	"github.com/mwm-io/gapi/examples/3-params/internal"
)

func main() {
	r := server.NewMux()

	server.AddHandlerFactory(r, "POST", "/body", internal.NewBodyHandler)
	server.AddHandlerFactory(r, "POST", "/body-with-params", internal.MakeBodyWithValidationHandler)

	server.AddHandlerFactory(r, "GET", "/path-params/{first}/{second}", internal.NewPathParamsHandler)
	server.AddHandlerFactory(r, "GET", "/query-params", internal.NewQueryParamsHandler)

	gLog.Info("Starting http server")
	if err := server.ServeAndHandleShutdown(r); err != nil {
		gLog.LogAny(err)
	}

	gLog.Info("Server stopped")
}
