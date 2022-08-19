package main

import (
	"github.com/mwm-io/gapi/examples/hello-world/internal"
	"github.com/mwm-io/gapi/request"
	"github.com/mwm-io/gapi/router"
)

func main() {
	r := router.Create()

	request.AddHandler(r, "GET", "/hello", internal.HelloWorldHandlerF())

	router.Handle(r)
}
