package main

import (
	"context"

	gLog "github.com/mwm-io/gapi/log"
	"github.com/mwm-io/gapi/openapi"
	"github.com/mwm-io/gapi/server"

	"github.com/mwm-io/gapi/examples/5-mcp/internal"
)

// This example shows how to expose your GAPI handlers as MCP tools.
// All documented handlers are automatically available as MCP tools via the /mcp endpoint.
//
// You can:
//   - Opt out a handler with builder.WithMCP(false)
//   - Override the auto-generated tool name with builder.WithMCPToolName("custom_name")
//
// To test:
//  1. Run this example: go run .
//  2. OpenAPI doc is available at http://localhost:8080/
//  3. Connect an MCP client (e.g. Claude Code, MCP Inspector) to http://localhost:8080/mcp
//  4. The client will see tools like: find_user, get_searchUsers, post_createUser, etc.
//  5. The /health endpoint is NOT exposed as an MCP tool (opted out)
func main() {
	ctx := context.Background()

	openapi.Config.DocPageTitle = "MCP Example | GAPI"

	r := server.NewMux()

	// CRUD handlers — all automatically exposed as MCP tools
	server.AddHandlerFactory(r, "GET", "/users", internal.SearchHandler)
	server.AddHandlerFactory(r, "POST", "/users", internal.CreateHandler)
	server.AddHandlerFactory(r, "GET", "/users/{id}", internal.GetOneHandler)
	server.AddHandlerFactory(r, "PUT", "/users/{id}", internal.UpdateHandler)
	server.AddHandlerFactory(r, "DELETE", "/users/{id}", internal.DeleteHandler)

	// Health check — opted out of MCP with WithMCP(false)
	server.AddHandler(r, "GET", "/health", internal.HealthHandler())

	// OpenAPI documentation at / and /openapi.json
	if err := server.AddDocHandlers(r); err != nil {
		gLog.Error(ctx).LogError(err)
	}

	// MCP server at /mcp — all documented handlers become MCP tools
	server.AddMCPServer(r)

	gLog.Info(ctx).LogMsg("Starting http server")
	gLog.Info(ctx).LogMsg("OpenAPI doc visualizer: http://localhost:8080/")
	gLog.Info(ctx).LogMsg("OpenAPI spec : http://localhost:8080/openapi.json")
	gLog.Info(ctx).LogMsg("MCP endpoint: http://localhost:8080/mcp")

	if err := server.ServeAndHandleShutdown(r); err != nil {
		gLog.Error(ctx).LogError(err)
	}

	gLog.Info(ctx).LogMsg("Server stopped")
}
