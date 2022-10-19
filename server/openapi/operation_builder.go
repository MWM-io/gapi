package openapi

import (
	"fmt"
	"strconv"

	"github.com/swaggest/openapi-go/openapi3"

	"github.com/mwm-io/gapi/errors"
)

// OperationBuilder is a builder to simplify the documentation of an operation
type OperationBuilder struct {
	operation  *openapi3.Operation
	reflector  *openapi3.Reflector
	httpMethod string
	path       string

	err []error
}

// NewOperationBuilder returns a new doc.OperationBuilder
func NewOperationBuilder(reflector *openapi3.Reflector, method, path string) *OperationBuilder {
	return &OperationBuilder{
		reflector:  reflector,
		operation:  new(openapi3.Operation),
		httpMethod: method,
		path:       path,
	}
}

// Operation return the openapi3.Operation used to compute doc for current operation
func (b OperationBuilder) Operation() *openapi3.Operation {
	return b.operation
}

// Reflector return the openapi3.Reflector used to compute doc for current operation
func (b OperationBuilder) Reflector() *openapi3.Reflector {
	return b.reflector
}

// Commit submit pending changes and return errors that was generated when building the operation
func (b *OperationBuilder) Commit() *OperationBuilder {
	if err := b.reflector.Spec.AddOperation(b.httpMethod, b.path, *b.operation); err != nil {
		b.err = append(b.err, err)
	}

	return b
}

// Error return the error that was generated when building the operation
func (b *OperationBuilder) Error() error {
	if b.err == nil {
		return nil
	}

	var err error
	for _, item := range b.err {
		err = errors.Err(fmt.Sprintf("%+v", item), err)
	}

	return err
}

// WithSummary set a Summary (Title) to the operation
func (b *OperationBuilder) WithSummary(summary string) *OperationBuilder {
	b.operation.WithSummary(summary)
	return b
}

// WithDescription set a description (additional explanation) to the operation
func (b *OperationBuilder) WithDescription(description string) *OperationBuilder {
	b.operation.WithDescription(description)
	return b
}

// WithTags set tags to the operation: used to organise you operation in sections
func (b *OperationBuilder) WithTags(tags ...string) *OperationBuilder {
	b.operation.WithTags(tags...)
	return b
}

// WithBody configure a request body to the operation
// Allowed options :
// - WithDescription to add a description to body
// - WithExample to add example(s) as body
// TODO: Find a way to support non json body like CSV, files, multi part, url encoded ...
func (b *OperationBuilder) WithBody(body interface{}, options ...BuilderOption) *OperationBuilder {
	var c config
	c.applyOptions(options)

	err := b.reflector.SetRequest(b.operation, body, b.httpMethod)
	if err != nil {
		b.err = append(b.err, err)
		return b
	}

	if c.description != "" {
		b.operation.RequestBodyEns().RequestBodyEns().Description = &c.description
	}

	if len(c.examples) != 0 {
		for mimeType, val := range b.operation.RequestBodyEns().RequestBodyEns().Content {
			for exampleKey, example := range c.examples {
				val.Examples[exampleKey] = openapi3.ExampleOrRef{
					Example: &example,
				}
			}

			b.operation.RequestBodyEns().RequestBodyEns().Content[mimeType] = val
		}
	}

	return b
}

// WithBodyExample set an example to request body to the operation
func (b *OperationBuilder) WithBodyExample(value interface{}) *OperationBuilder {
	for mimeType, val := range b.operation.RequestBodyEns().RequestBodyEns().Content {
		val.WithExample(value)
		b.operation.RequestBodyEns().RequestBodyEns().Content[mimeType] = val
	}

	return b
}

// WithParams configure path and query parameters to the operation
// To set path parameters use a struct with path tag
// To set query parameters use a struct with query tag
func (b *OperationBuilder) WithParams(body interface{}) *OperationBuilder {
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
func (b *OperationBuilder) WithResponse(output interface{}, options ...BuilderOption) *OperationBuilder {
	c := config{
		description: "",
		examples:    nil,
		mimeType:    "application/json",
		statusCode:  200,
	}

	if output == nil {
		c.statusCode = 204
	}

	c.applyOptions(options)

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

	if len(c.examples) != 0 {
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
