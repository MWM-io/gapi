package openapi

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/mwm-io/gapi/handler"
)

// RapiDocHandler is a server.SpecOpenAPIHandler that will serve a html page with rapidoc loading the given openAPIJsonURL
type RapiDocHandler struct {
	handler.WithMiddlewares
	openAPIJsonURL string
}

// NewRapiDocHandler build a new RapiDocHandler.
func NewRapiDocHandler(middlewares ...handler.Middleware) handler.Handler {
	return &RapiDocHandler{
		openAPIJsonURL: Config.GetSpecOpenAPIURI(),
		WithMiddlewares: handler.WithMiddlewares{
			MiddlewareList: middlewares,
		},
	}
}

// Serve implements the handler.Handler interface
func (h RapiDocHandler) Serve(_ http.ResponseWriter, _ *http.Request) (interface{}, error) {
	// TODO add customisation options : like icon, theme etc...

	var favicon string
	if faviconURL := Config.GetFaviconURL(); faviconURL != "" {
		favicon = fmt.Sprintf(`<link rel="icon" type="image/x-icon" href="%s">`, faviconURL)
	}

	return strings.NewReader(`
<!doctype html> 
<html>
<style>
    rapi-doc::part(section-navbar-item section-navbar-tag) {
      color: var(--primary-color);
    }
</style>
<head>
  <title>` + Config.GetDocPageTitle() + `</title>
  ` + favicon + `
  <meta charset="utf-8"> <!-- Important: rapi-doc uses utf8 characters -->
  <script type="module" src="https://unpkg.com/rapidoc/dist/rapidoc-min.js"></script>
</head>
<body>
  <rapi-doc
    spec-url="` + h.openAPIJsonURL + `"
	theme="dark"
    primary-color="#f54c47"
    bg-color="#222c3d"
    text-color="#fff"
    show-header="false"
	render-style="read"
	schema-style="table"
	show-header="false"
  > </rapi-doc>
</body>
</html>`), nil
}

// RapiDocReceiverHandler is a server.SpecOpenAPIHandler that will serve the rapidoc oauth receiver
type RapiDocReceiverHandler struct {
	handler.WithMiddlewares
}

// NewRapiDocReceiverHandler build a new RapiDocReceiverHandler.
func NewRapiDocReceiverHandler(middlewares ...handler.Middleware) handler.Handler {
	return &RapiDocReceiverHandler{
		WithMiddlewares: handler.WithMiddlewares{
			MiddlewareList: middlewares,
		},
	}
}

// Serve implements the handler.Handler interface
func (h RapiDocReceiverHandler) Serve(_ http.ResponseWriter, _ *http.Request) (interface{}, error) {
	return strings.NewReader(`
<!doctype html>
<head>
  <script type="module" src="https://unpkg.com/rapidoc/dist/rapidoc-min.js"></script>
</head>

<body>
  <oauth-receiver> </oauth-receiver>
</body>
`), nil
}
