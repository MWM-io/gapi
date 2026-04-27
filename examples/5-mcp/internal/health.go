package internal

import (
	"net/http"

	"github.com/mwm-io/gapi/handler"
	"github.com/mwm-io/gapi/openapi"
)

type healthHandler struct{}

// HealthHandler returns a simple health check handler.
// This handler is documented for OpenAPI but opted out of MCP with WithMCP(false).
func HealthHandler() handler.Handler {
	return healthHandler{}
}

// Doc implements openapi.Documented.
// WithMCP(false) excludes this handler from MCP tool exposure — it will appear
// in the OpenAPI spec but not as an MCP tool.
func (h healthHandler) Doc(builder *openapi.DocBuilder) error {
	builder.
		WithSummary("Health check").
		WithDescription("Returns server health status. Not exposed as MCP tool.").
		WithTags("Internal").
		WithMCP(false).
		WithResponse(map[string]string{"status": "ok"})
	return nil
}

// Serve implements handler.Handler
func (h healthHandler) Serve(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return map[string]string{"status": "ok"}, nil
}
