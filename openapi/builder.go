package openapi

import (
	"strconv"

	"github.com/swaggest/openapi-go/openapi3"

	"github.com/mwm-io/gapi/errors"
	"github.com/mwm-io/gapi/utils"
)

// DocBuilder is a builder to simplify the documentation of an operation
type DocBuilder struct {
	operation  *openapi3.Operation
	reflector  *openapi3.Reflector
	httpMethod string
	path       string

	err []error

	// MCP fields
	mcpEnabled  *bool         // nil = default (true), explicit true/false overrides
	mcpToolName string        // custom tool name, empty = auto from operationID
	bodyPtr     interface{}   // raw body struct pointer
	paramPtrs   []interface{} // raw param struct pointers (path + query)
	responsePtr interface{}   // raw response struct pointer
}

// NewDocBuilder returns a new doc.DocBuilder
func NewDocBuilder(reflector *openapi3.Reflector, method, path string) *DocBuilder {
	return &DocBuilder{
		reflector:  reflector,
		operation:  new(openapi3.Operation),
		httpMethod: method,
		path:       path,
	}
}

// Operation return the openapi3.Operation used to compute doc for current operation
func (b *DocBuilder) Operation() *openapi3.Operation {
	return b.operation
}

// Reflector return the openapi3.Reflector used to builds OpenAPI Schema with reflected structures for current operation.
func (b *DocBuilder) Reflector() *openapi3.Reflector {
	return b.reflector
}

// Commit submit pending changes and return errors that was generated when building the operation
func (b *DocBuilder) Commit() *DocBuilder {
	if err := b.reflector.Spec.AddOperation(b.httpMethod, b.path, *b.operation); err != nil {
		b.err = append(b.err, err)
	}

	return b
}

// Error return the error that was generated when building the operation
func (b *DocBuilder) Error() error {
	if b.err == nil {
		return nil
	}

	var err error
	for _, item := range b.err {
		err = errors.Wrap(item)
	}

	return err
}

// WithSummary set a Summary (Title) to the operation
func (b *DocBuilder) WithSummary(summary string) *DocBuilder {
	b.operation.WithSummary(summary)
	defaultOperationID := utils.GenerateOperationID("/" + b.httpMethod + b.path)
	b.operation.ID = &defaultOperationID
	return b
}

// WithDescription set a description (additional explanation) to the operation
func (b *DocBuilder) WithDescription(description string) *DocBuilder {
	b.operation.WithDescription(description)
	return b
}

// WithTags set tags to the operation: used to organise you operation in sections
func (b *DocBuilder) WithTags(tags ...string) *DocBuilder {
	b.operation.WithTags(tags...)
	return b
}

// WithMCP sets whether this operation should be exposed as an MCP tool.
// Default is true for all documented operations.
func (b *DocBuilder) WithMCP(enabled bool) *DocBuilder {
	b.mcpEnabled = &enabled
	return b
}

// WithMCPToolName overrides the auto-generated MCP tool name.
func (b *DocBuilder) WithMCPToolName(name string) *DocBuilder {
	b.mcpToolName = name
	return b
}

// IsMCPEnabled returns whether this operation should be exposed as an MCP tool.
// Returns true by default unless explicitly disabled.
func (b *DocBuilder) IsMCPEnabled() bool {
	if b.mcpEnabled == nil {
		return true
	}
	return *b.mcpEnabled
}

// MCPToolName returns the custom MCP tool name, or empty string if not set.
func (b *DocBuilder) MCPToolName() string {
	return b.mcpToolName
}

// BodyPtr returns the raw body struct pointer stored during WithBody.
func (b *DocBuilder) BodyPtr() interface{} {
	return b.bodyPtr
}

// ParamPtrs returns the raw parameter struct pointers stored during WithParams.
func (b *DocBuilder) ParamPtrs() []interface{} {
	return b.paramPtrs
}

// ResponsePtr returns the raw response struct pointer stored during WithResponse.
func (b *DocBuilder) ResponsePtr() interface{} {
	return b.responsePtr
}

// HTTPMethod returns the HTTP method for this operation.
func (b *DocBuilder) HTTPMethod() string {
	return b.httpMethod
}

// Path returns the URL path for this operation.
func (b *DocBuilder) Path() string {
	return b.path
}

// OperationID returns the operation ID, or empty string if not set.
func (b *DocBuilder) OperationID() string {
	if b.operation.ID != nil {
		return *b.operation.ID
	}
	return ""
}

// Summary returns the operation summary, or empty string if not set.
func (b *DocBuilder) Summary() string {
	if b.operation.Summary != nil {
		return *b.operation.Summary
	}
	return ""
}

// Description returns the operation description, or empty string if not set.
func (b *DocBuilder) Description() string {
	if b.operation.Description != nil {
		return *b.operation.Description
	}
	return ""
}

// WithBody configure a request body to the operation
// Allowed options :
// - WithDescription to add a description to body
// - WithExample to add example(s) as body
// TODO: Find a way to support non json body like CSV, files, multi part, url encoded ...
func (b *DocBuilder) WithBody(body interface{}, options ...BuilderOption) *DocBuilder {
	b.bodyPtr = body

	var c builderOptions
	c.applyOptions(options...)

	err := b.reflector.SetRequest(b.operation, body, b.httpMethod)
	if err != nil {
		b.err = append(b.err, err)
		return b
	}

	if c.description != "" {
		b.operation.RequestBodyEns().RequestBodyEns().Description = &c.description
	}

	if len(c.examples) == 0 {
		c.examples = make(map[string]openapi3.Example)
		exampleName := "default"
		c.examples[exampleName] = openapi3.Example{
			Summary: &exampleName,
			Value:   &body,
		}
	}

	for mimeType, val := range b.operation.RequestBodyEns().RequestBodyEns().Content {
		for exampleKey, example := range c.examples {
			val.WithExamplesItem(exampleKey, openapi3.ExampleOrRef{
				Example: &example,
			})
		}

		b.operation.RequestBodyEns().RequestBodyEns().WithContentItem(mimeType, val)
	}

	return b
}

// WithBodyExample set an example to request body to the operation
func (b *DocBuilder) WithBodyExample(value interface{}) *DocBuilder {
	for mimeType, val := range b.operation.RequestBodyEns().RequestBodyEns().Content {
		b.operation.RequestBodyEns().RequestBodyEns().WithContentItem(mimeType, *val.WithExample(value))
	}

	return b
}

// WithParams configure path and query parameters to the operation
// To set path parameters use a struct with 'path' tag
// To set query parameters use a struct with 'query' tag
func (b *DocBuilder) WithParams(body interface{}) *DocBuilder {
	b.paramPtrs = append(b.paramPtrs, body)

	if err := b.reflector.SetRequest(b.operation, body, b.httpMethod); err != nil {
		b.err = append(b.err, err)
	}

	return b
}

// WithResponse configure a response for current operation
// Allowed options :
// - WithDescription to add a description to response
// - WithExample to add example(s) as response
// - WithMimeType to set a custom contentType (default to json)
// - WithStatusCode to set a specific status code. Default value are 204 for nil value and 200 for non nil value
func (b *DocBuilder) WithResponse(output interface{}, options ...BuilderOption) *DocBuilder {
	b.responsePtr = output

	c := builderOptions{
		description: "",
		examples:    nil,
		mimeType:    "application/json",
		statusCode:  200,
	}

	if output == nil {
		c.statusCode = 204
		c.mimeType = ""
	}

	c.applyOptions(options...)

	err := b.reflector.SetupResponse(openapi3.OperationContext{
		Operation:         b.operation,
		Output:            output,
		HTTPStatus:        c.statusCode,
		RespContentType:   c.mimeType,
		RespHeaderMapping: c.headers,
	})
	if err != nil {
		b.err = append(b.err, err)
	}

	statusCodeStr := strconv.Itoa(c.statusCode)
	resp := b.operation.Responses.MapOfResponseOrRefValues[statusCodeStr]
	if c.description != "" {
		resp.ResponseEns().WithDescription(c.description)
	}

	if len(c.examples) == 0 && output != nil {
		c.examples = make(map[string]openapi3.Example)
		exampleName := "default"
		c.examples[exampleName] = openapi3.Example{
			Summary:     &exampleName,
			Description: &exampleName,
			Value:       &output,
		}
	}

	if len(c.examples) != 0 && output != nil {
		contentResp := resp.ResponseEns().Content[c.mimeType]

		for key, item := range c.examples {
			contentResp.WithExamplesItem(key, openapi3.ExampleOrRef{
				Example: &item,
			})
		}

		resp.ResponseEns().WithContentItem(c.mimeType, contentResp)
	}

	b.operation.Responses.WithMapOfResponseOrRefValuesItem(statusCodeStr, resp)

	return b
}

// WithError configure an error for current operation
// Allowed options :
// - WithDescription to add a description to error response
func (b *DocBuilder) WithError(statusCode int, kind, message string, options ...BuilderOption) *DocBuilder {
	c := builderOptions{
		examples:    nil,
		statusCode:  statusCode,
		description: message,
		mimeType:    "application/json",
	}

	c.applyOptions(options...)

	exampleValue := errors.HttpError{
		Message: message,
		Kind:    kind,
	}
	err := b.reflector.SetupResponse(openapi3.OperationContext{
		Operation:         b.operation,
		Output:            exampleValue,
		HTTPStatus:        c.statusCode,
		RespContentType:   c.mimeType,
		RespHeaderMapping: c.headers,
	})
	if err != nil {
		b.err = append(b.err, err)
	}

	statusCodeStr := strconv.Itoa(c.statusCode)
	resp := b.operation.Responses.MapOfResponseOrRefValues[statusCodeStr]
	if c.description != "" {
		resp.ResponseEns().WithDescription(c.description)
	}

	// Use given data as example : set an example
	jsonResp := resp.ResponseEns().Content[c.mimeType]
	jsonResp.WithExamplesItem(kind, openapi3.ExampleOrRef{
		Example: new(openapi3.Example).WithValue(exampleValue),
	})
	resp.ResponseEns().WithContentItem(c.mimeType, jsonResp)

	b.operation.Responses.WithMapOfResponseOrRefValuesItem(statusCodeStr, resp)

	return b
}

// WithOperationID set an operationID to the operation
func (b *DocBuilder) WithOperationID(operationID string) *DocBuilder {
	b.operation.ID = &operationID
	return b
}
