package server

import (
	"github.com/gorilla/mux"

	"github.com/mwm-io/gapi/mcp"
	"github.com/mwm-io/gapi/openapi"
)

// AddMCPServer registers an MCP server endpoint on the router.
// It exposes all documented handlers as MCP tools via Streamable HTTP transport.
// The MCP path is automatically added to openapi.Config.IgnoredPaths so it
// doesn't appear in the OpenAPI spec.
func AddMCPServer(r *mux.Router) {
	mcpPath := mcp.MCPConfig.GetPath()

	// Exclude the MCP endpoint from OpenAPI documentation
	openapi.Config.IgnoredPaths = append(openapi.Config.IgnoredPaths, mcpPath)

	mcpHandler := mcp.NewMCPHandler(r, mcp.MCPConfig)
	r.PathPrefix(mcpPath).Handler(mcpHandler)
}
