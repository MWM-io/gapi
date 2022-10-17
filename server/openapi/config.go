package openapi

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/swaggest/openapi-go/openapi3"

	"github.com/mwm-io/gapi/server"
)

const (
	defaultDocURI     = "/"
	defaultDocJSONURI = "/openapi.json"
)

type Config struct {
	DocURI       string
	DocJSONURI   string
	IgnoredPaths []string
	Auth         Authorization
	Reflector    *openapi3.Reflector
}

func (c Config) docURI() string {
	if c.DocURI != "" {
		return c.DocURI
	}

	return defaultDocURI
}

func (c Config) docJSONURI() string {
	if c.DocJSONURI != "" {
		return c.DocJSONURI
	}

	return defaultDocJSONURI
}

func (c Config) docReceiverURI() string {
	return fmt.Sprintf("%s%s", c.docURI(), "oauth-receiver.html")
}

func (c Config) ignoredPaths() []string {
	ignoredPaths := []string{
		c.docURI(),
		c.docJSONURI(),
		c.docReceiverURI(),
	}

	return append(ignoredPaths, c.IgnoredPaths...)
}

func (c Config) reflector() *openapi3.Reflector {
	if c.Reflector != nil {
		return c.Reflector
	}

	reflector := &openapi3.Reflector{}
	reflector.SpecEns()

	return reflector
}

// AddRapidocHandlers will add the necessary handlers to serve a rapidoc endpoint.
func AddRapidocHandlers(r *mux.Router, config Config) error {
	reflector := config.reflector()
	err := PopulateReflector(reflector, r, config.ignoredPaths())
	if err != nil {
		return err
	}

	server.AddHandler(r, http.MethodGet, config.docURI(), NewRapiDocHandler(config.docJSONURI(), config.Auth))
	server.AddHandler(r, http.MethodGet, config.docReceiverURI(), RapiDocReceiverHandler)
	server.AddHandler(r, http.MethodGet, config.docJSONURI(), NewHandler(reflector, config.Auth))

	return nil
}
