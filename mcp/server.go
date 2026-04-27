package mcp

import (
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/mwm-io/gapi/openapi"
)

// NewMCPHandler creates an http.Handler that serves the MCP protocol over
// Streamable HTTP transport. It lazily discovers tools from the router on
// first connection, similar to how OpenAPI spec generation works.
func NewMCPHandler(router *mux.Router, config Config) http.Handler {
	mcpServer := mcpsdk.NewServer(
		&mcpsdk.Implementation{
			Name:    config.GetServerName(openapi.Config.GetDocPageTitle()),
			Version: config.GetServerVersion(),
		},
		nil,
	)

	var once sync.Once
	var populateErr error

	// Wrap the MCP HTTP handler to lazily populate tools on first request
	streamHandler := mcpsdk.NewStreamableHTTPHandler(
		func(r *http.Request) *mcpsdk.Server {
			once.Do(func() {
				populateErr = PopulateMCPServer(mcpServer, router, config)
			})
			if populateErr != nil {
				return nil
			}
			return mcpServer
		},
		nil,
	)

	return streamHandler
}
