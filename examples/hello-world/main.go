package main

import (
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/mwm-io/gapi/errors"
	"github.com/mwm-io/gapi/middleware"

	gLog "github.com/mwm-io/gapi/log"
	"github.com/mwm-io/gapi/server"
)

func main() {
	r := server.NewMux()

	// If you don't add any middlewares, handlers will work as native http handlers.
	server.AddHandler(r, "GET", "/hello-world", server.HandlerF(func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
		w.Write([]byte("hello-world"))
		return nil, nil
	}))

	// You can add middlewares to use the returned values (interface{}, error) to return your response.
	server.AddHandler(r, "GET", "/hello-world-1", server.HandlerF(func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
		return "hello-world", nil
	}, middleware.ResponseWriterMiddleware{
		Marshalers: map[string]middleware.Marshaler{"text/plain": middleware.MarshalerFunc(func(v interface{}) ([]byte, error) {
			return []byte(fmt.Sprintf("%#v", v)), nil
		})},
		DefaultContentType: "text/plain",
		Response:           "response",
	}))

	// Use the core to add all basic middlewares needed, with a predefined configuration.
	server.AddHandler(r, "GET", "/hello-world-core", server.HandlerF(func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
		return "hello-world", nil
	}, middleware.Core().Middlewares()...))

	// You can also return an object, that will be serialized by the response_writer middleware.
	server.AddHandler(r, "GET", "/hello-world-serialized", server.HandlerF(func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
		return struct {
			XMLName xml.Name `xml:"response"`
			Title   string   `xml:"title" json:"title"`
		}{Title: "hello-world"}, nil
	}, middleware.Core().Middlewares()...))

	// You can also add option to custom the core's middlewares. Here we will always return a json response.
	server.AddHandler(r, "GET", "/hello-world-json", server.HandlerF(func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
		return struct {
			Title string `json:"title"`
		}{Title: "hello-world"}, nil
	}, middleware.Core(middleware.WithForcedContentType("application/json")).Middlewares()...))

	// The response write in the core will also read and return the returned error.
	server.AddHandler(r, "GET", "/hello-world-error", server.HandlerF(func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
		return nil, errors.Err("my error")
	}, middleware.Core().Middlewares()...))

	gLog.Info("Starting http server")

	if err := server.ServeAndHandleShutdown(r); err != nil {
		gLog.LogAny(err)
	}

	gLog.Info("Server stopped")
}
