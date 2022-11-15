package main

import (
	"github.com/mwm-io/gapi/examples/2-encoders/internal/err"
	"github.com/mwm-io/gapi/examples/2-encoders/internal/hello"
	gLog "github.com/mwm-io/gapi/log"
	"github.com/mwm-io/gapi/server"
)

func main() {
	r := server.NewMux()

	server.AddHandler(r, "GET", "/hello/json", hello.MakeJSONResponseHandler())
	server.AddHandler(r, "GET", "/hello/xml", hello.MakeXMLResponseHandler())
	server.AddHandler(r, "GET", "/hello/auto", hello.MakeAutoResponseHandler())

	server.AddHandler(r, "GET", "/error/json", err.MakeJSONResponseHandler())
	server.AddHandler(r, "GET", "/error/xml", err.MakeXMLResponseHandler())
	server.AddHandler(r, "GET", "/error/auto", err.MakeAutoResponseHandler())

	gLog.Info("Starting http server")
	if errServe := server.ServeAndHandleShutdown(r); errServe != nil {
		gLog.LogAny(errServe)
	}

	gLog.Info("Server stopped")
}
