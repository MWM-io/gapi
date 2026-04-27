package mcp

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/gorilla/mux"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/swaggest/openapi-go/openapi3"

	"github.com/mwm-io/gapi/openapi"
)

// PopulateMCPServer walks the mux router and registers each documented handler
// as an MCP tool on the given MCP server.
func PopulateMCPServer(mcpServer *mcpsdk.Server, router *mux.Router, config Config) error {
	// Create a throwaway reflector for doc generation — we don't care about the OpenAPI spec here,
	// we just need the DocBuilder metadata.
	reflector := new(openapi3.Reflector)
	reflector.SpecEns().Info.WithTitle(config.GetServerName(""))

	ignoredPaths := allIgnoredPaths(config)

	return router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, errPath := route.GetPathTemplate()
		if errPath != nil {
			return nil
		}

		for _, ignoredPath := range ignoredPaths {
			if path == ignoredPath {
				return nil
			}
		}

		methods, errMethods := route.GetMethods()
		if errMethods != nil {
			return nil
		}

		for _, method := range methods {
			httpHandler := route.GetHandler()

			builder, err := openapi.BuildDocBuilder(reflector, httpHandler, method, path)
			if err != nil {
				return err
			}

			if !builder.IsMCPEnabled() {
				continue
			}

			// Generate tool name
			toolName := builder.MCPToolName()
			if toolName == "" {
				toolName = generateToolName(method, builder.OperationID(), path)
			}

			// Generate input schema
			schemaJSON, paramMap, err := GenerateInputSchema(builder.ParamPtrs(), builder.BodyPtr())
			if err != nil {
				return err
			}

			// Parse the schema into a map for the MCP SDK
			var schemaMap map[string]interface{}
			if err := json.Unmarshal(schemaJSON, &schemaMap); err != nil {
				return err
			}

			// Build tool description
			description := builder.Summary()
			if builder.Description() != "" {
				if description != "" {
					description += " — "
				}
				description += builder.Description()
			}

			// Create the tool route for execution
			tr := &toolRoute{
				method:      method,
				pathTpl:     path,
				httpHandler: httpHandler,
				paramMap:    paramMap,
			}

			// Register the tool
			mcpServer.AddTool(
				&mcpsdk.Tool{
					Name:        toolName,
					Description: description,
					InputSchema: schemaMap,
				},
				func(ctx context.Context, req *mcpsdk.CallToolRequest) (*mcpsdk.CallToolResult, error) {
					return tr.Execute(ctx, req)
				},
			)
		}

		return nil
	})
}

// generateToolName creates a tool name from the HTTP method and operationID or path.
func generateToolName(method, operationID, path string) string {
	method = strings.ToLower(method)

	if operationID != "" {
		return method + "_" + operationID
	}

	// Fallback: build from path
	name := strings.ReplaceAll(path, "/", "_")
	name = strings.ReplaceAll(name, "{", "")
	name = strings.ReplaceAll(name, "}", "")
	name = strings.Trim(name, "_")

	return method + "_" + name
}

// allIgnoredPaths combines MCP-specific ignored paths with OpenAPI doc/spec paths.
func allIgnoredPaths(config Config) []string {
	paths := []string{config.GetPath()}
	paths = append(paths, config.IgnoredPaths...)

	// Include the OpenAPI documentation endpoints so they don't become MCP tools
	paths = append(paths,
		openapi.Config.GetDocURI(),
		openapi.Config.GetSpecOpenAPIURI(),
		openapi.Config.GetAuthReceiverURI(),
	)
	paths = append(paths, openapi.Config.IgnoredPaths...)

	return paths
}
