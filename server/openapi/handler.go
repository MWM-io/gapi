package openapi

import (
	"encoding/json"
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

// AddRapidocHandlers will add the necessary handlers to serve a rapidoc endpoint:
// 2 endpoints to serve rapidoc.html and oauth-receiver.html from rapidoc
// and one endpoint to serve the json openapi definition of your API.
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

// Handler is a server.Handler that will return a json openapi definition of the given reflector.
type Handler struct {
	reflector *openapi3.Reflector
	auth      Authorization
}

// NewHandler builds a new Handler, serving the api definition from the openapi3.Reflector,
// and checking auth access with the given Authorization.
func NewHandler(reflector *openapi3.Reflector, auth Authorization) *Handler {
	return &Handler{
		reflector: reflector,
		auth:      auth,
	}
}

// Serve implements the server.Handler interface
func (h *Handler) Serve(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	if h.auth != nil {
		authorized, err := h.auth.Authorize(w, r)
		if err != nil {
			return nil, err
		}
		if !authorized {
			return h.auth.Login(w, r)
		}
	}

	body, err := json.Marshal(h.reflector.Spec)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	_, err = w.Write(body)
	return nil, err
}

// Config contains all the options for your API documentation.
type Config struct {
	// DocURI is the URI for the display of your documentation.
	DocURI string
	// DocJSONURI is the URL for the json openapi definition of your API.
	DocJSONURI string
	// IgnoredPaths are the paths that shouldn't be included in the documentation.
	IgnoredPaths []string
	// Auth is your auth system to protect your documentation.
	Auth Authorization
	// Reflector is the base documentation for your API,
	// so you can add generic information before automatically adding your handlers' documentation.
	Reflector *openapi3.Reflector
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
