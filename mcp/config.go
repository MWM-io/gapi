package mcp

// DefaultMCPPath is the default HTTP path for the MCP server endpoint.
const DefaultMCPPath = "/mcp"

// MCPConfig is the default MCP server config. You can update all fields to change:
//
// - ServerName: the name of the MCP server.
// - ServerVersion: the version of the MCP server.
// - Path: the HTTP path for the MCP endpoint.
// - IgnoredPaths: the paths that shouldn't be exposed as MCP tools.
var MCPConfig Config

// Config contains the configuration for the MCP server.
type Config struct {
	// ServerName is the name of the MCP server. Defaults to openapi.Config.DocPageTitle.
	ServerName string
	// ServerVersion is the version of the MCP server. Defaults to "1.0.0".
	ServerVersion string
	// Path is the HTTP path for the Streamable HTTP transport. Defaults to "/mcp".
	Path string
	// IgnoredPaths are the routes that shouldn't be exposed as MCP tools.
	IgnoredPaths []string
}

// GetPath returns the configured MCP path, or the default if not set.
func (c Config) GetPath() string {
	if c.Path != "" {
		return c.Path
	}
	return DefaultMCPPath
}

// GetServerName returns the configured server name, or the fallback.
func (c Config) GetServerName(fallback string) string {
	if c.ServerName != "" {
		return c.ServerName
	}
	if fallback != "" {
		return fallback
	}
	return "gapi-mcp-server"
}

// GetServerVersion returns the configured server version, or "1.0.0".
func (c Config) GetServerVersion() string {
	if c.ServerVersion != "" {
		return c.ServerVersion
	}
	return "1.0.0"
}
