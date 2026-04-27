package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gorilla/mux"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// toolRoute holds the information needed to execute a GAPI handler from an MCP tool call.
type toolRoute struct {
	method      string
	pathTpl     string            // e.g., "/users/{id}"
	httpHandler http.Handler      // the defaultHandleEngine
	paramMap    map[string]ParamSource
}

// Execute runs the GAPI handler for this tool route, decomposing the flat MCP arguments
// back into path vars, query params, and body fields, then returning the result.
func (t *toolRoute) Execute(ctx context.Context, req *mcpsdk.CallToolRequest) (*mcpsdk.CallToolResult, error) {
	// Parse the arguments from the request
	var arguments map[string]interface{}
	if req.Params.Arguments != nil {
		if err := json.Unmarshal(req.Params.Arguments, &arguments); err != nil {
			return newErrorResult(fmt.Sprintf("failed to parse arguments: %v", err)), nil
		}
	}
	if arguments == nil {
		arguments = make(map[string]interface{})
	}

	// Decompose flat arguments into path, query, and body
	pathVars := make(map[string]string)
	queryVars := make(map[string]string)
	bodyFields := make(map[string]interface{})

	for key, value := range arguments {
		source, ok := t.paramMap[key]
		if !ok {
			// Unknown argument — put it in the body as a best-effort fallback
			bodyFields[key] = value
			continue
		}

		switch source.Source {
		case "path":
			pathVars[key] = fmt.Sprintf("%v", value)
		case "query":
			queryVars[key] = fmt.Sprintf("%v", value)
		case "body":
			bodyFields[key] = value
		}
	}

	// Build URL by filling path template with path vars
	url := t.pathTpl
	for key, val := range pathVars {
		url = strings.ReplaceAll(url, "{"+key+"}", val)
	}

	// Build query string
	if len(queryVars) > 0 {
		parts := make([]string, 0, len(queryVars))
		for k, v := range queryVars {
			parts = append(parts, k+"="+v)
		}
		url += "?" + strings.Join(parts, "&")
	}

	// Build request body
	var bodyReader *bytes.Reader
	if len(bodyFields) > 0 {
		bodyBytes, err := json.Marshal(bodyFields)
		if err != nil {
			return newErrorResult(fmt.Sprintf("failed to encode body: %v", err)), nil
		}
		bodyReader = bytes.NewReader(bodyBytes)
	} else {
		bodyReader = bytes.NewReader(nil)
	}

	// Create the HTTP request
	httpReq := httptest.NewRequest(t.method, url, bodyReader)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	// Propagate context
	httpReq = httpReq.WithContext(ctx)

	// Set mux path vars so PathParameters middleware works
	if len(pathVars) > 0 {
		httpReq = mux.SetURLVars(httpReq, pathVars)
	}

	// Execute the full GAPI pipeline
	recorder := httptest.NewRecorder()
	t.httpHandler.ServeHTTP(recorder, httpReq)

	// Format the response
	responseBody := recorder.Body.String()
	statusCode := recorder.Code

	if statusCode >= 200 && statusCode < 300 {
		return &mcpsdk.CallToolResult{
			Content: []mcpsdk.Content{
				&mcpsdk.TextContent{Text: responseBody},
			},
		}, nil
	}

	// Error response
	return &mcpsdk.CallToolResult{
		Content: []mcpsdk.Content{
			&mcpsdk.TextContent{Text: responseBody},
		},
		IsError: true,
	}, nil
}

// newErrorResult creates a CallToolResult with IsError=true.
func newErrorResult(message string) *mcpsdk.CallToolResult {
	return &mcpsdk.CallToolResult{
		Content: []mcpsdk.Content{
			&mcpsdk.TextContent{Text: message},
		},
		IsError: true,
	}
}
