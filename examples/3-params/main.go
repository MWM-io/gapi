package main

import (
	"github.com/mwm-io/gapi/examples/3-params/internal"
	gLog "github.com/mwm-io/gapi/log"
	"github.com/mwm-io/gapi/server"
)

func main() {
	r := server.NewMux()

	server.AddHandlerFactory(r, "POST", "/body", internal.NewBodyHandler)
	server.AddHandlerFactory(r, "GET", "/path-params/{first}/{second}", internal.NewPathParamsHandler)
	server.AddHandlerFactory(r, "GET", "/query-params", internal.NewQueryParamsHandler)

	gLog.Info("Starting http server")
	if err := server.ServeAndHandleShutdown(r); err != nil {
		gLog.LogAny(err)
	}

	gLog.Info("Server stopped")
}
