package main

import (
	"github.com/mwm-io/gapi/examples/hello-world/internal"
	"github.com/mwm-io/gapi/request"
	"github.com/mwm-io/gapi/router"
)

func main() {
	r := router.Create()

	request.AddHandler(r, "GET", "/json/hello", internal.JsonHelloWorldHandlerF())
	request.AddHandler(r, "GET", "/xml/hello", internal.XmlHelloWorldHandlerF())
	request.AddHandler(r, "GET", "/error/hello", internal.ErrorHelloWorldHandlerF())

	router.Handle(r)
}
