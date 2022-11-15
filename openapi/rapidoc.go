package openapi

import (
	"net/http"
	"strings"
)

// RapiDocHandler is a server.SpecOpenAPIHandler that will serve a html page with rapidoc loading the given openAPIJsonURL
type RapiDocHandler struct {
	openAPIJsonURL string
}

// NewRapiDocHandler build a new RapiDocHandler.
func NewRapiDocHandler() RapiDocHandler {
	return RapiDocHandler{
		openAPIJsonURL: Config.GetSpecOpenAPIURI(),
	}
}

// Serve implements the server.SpecOpenAPIHandler interface
func (h RapiDocHandler) Serve(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	// TODO add customisation options : like icon, theme etc...
	// TODO add auth middleware as option
	return strings.NewReader(`
<!doctype html> 
<html>
<head>
  <meta charset="utf-8"> <!-- Important: rapi-doc uses utf8 characters -->
  <script type="module" src="https://unpkg.com/rapidoc/dist/rapidoc-min.js"></script>
</head>
<body>
  <rapi-doc
    spec-url="` + h.openAPIJsonURL + `"
	bg-color="#14191f"
	text-color="#aec2e0"
	theme="dark"
	render-style="read"
	schema-style="table"
  > </rapi-doc>
</body>
</html>`), nil
}

// RapiDocReceiverHandler is a server.SpecOpenAPIHandler that will serve the rapidoc oauth receiver
func RapiDocReceiverHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
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
