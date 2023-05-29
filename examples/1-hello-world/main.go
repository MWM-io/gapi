package main

import (
	"context"
	"net/http"

	"github.com/mwm-io/gapi/handler"
	gLog "github.com/mwm-io/gapi/log"
	"github.com/mwm-io/gapi/server"
)

func main() {
	ctx := context.Background()
	r := server.NewMux()

	server.AddHandler(r, "GET", "/", handler.Func(HelloWorldHandler))

	gLog.Info(ctx, "Starting http server")
	if err := server.ServeAndHandleShutdown(r); err != nil {
		gLog.Emergency(ctx, err.Error())
	}

	gLog.Info(ctx, "Server stopped")
}

// HelloWorldHandler is the simplest handler with core middlewares.
// Reply "Hello World" string marshaled using Accept header
func HelloWorldHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return "Hello world", nil
}
