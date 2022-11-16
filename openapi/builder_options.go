package openapi

import (
	"github.com/swaggest/openapi-go/openapi3"
)

type builderOptions struct {
	summary     string
	description string
	examples    map[string]openapi3.Example
	mimeType    string
	statusCode  int
	headers     map[string]string
}

func (c *builderOptions) applyOptions(options ...BuilderOption) {
	for _, option := range options {
		option(c)
	}
}

// BuilderOption can be used to override default value or add precision to documented element
type BuilderOption func(c *builderOptions)

// WithDescription adds a description to documented element
func WithDescription(description string) BuilderOption {
	return func(c *builderOptions) {
		c.description = description
	}
}

// WithExample add an example to documented element
func WithExample(exampleName string, value interface{}, options ...BuilderOption) BuilderOption {
	return func(c *builderOptions) {
		var exampleOptions builderOptions
		exampleOptions.applyOptions(options...)

		var description *string
		if exampleOptions.description != "" {
			description = &exampleOptions.description
		}

		if c.examples == nil {
			c.examples = make(map[string]openapi3.Example)
		}

		c.examples[exampleName] = openapi3.Example{
			Summary:     &exampleName,
			Description: description,
			Value:       &value,
		}
	}
}

// WithMimeType add precision about type/format (or override default value) to documented element
func WithMimeType(valueType string) BuilderOption {
	return func(c *builderOptions) {
		c.mimeType = valueType
	}
}

// WithStatusCode override default status code to documented element
func WithStatusCode(statusCode int) BuilderOption {
	return func(c *builderOptions) {
		c.statusCode = statusCode
	}
}

// WithRedirect add a redirection response to the documented element.
func WithRedirect(statusCode int, location string) BuilderOption {
	return func(c *builderOptions) {
		WithStatusCode(statusCode)(c)

		WithHeaders(map[string]string{
			"Location": location,
		})(c)
	}
}

// WithHeaders add headers the documented element.
func WithHeaders(headers map[string]string) BuilderOption {
	return func(c *builderOptions) {
		if c.headers == nil {
			c.headers = make(map[string]string)
		}

		for k, v := range c.headers {
			c.headers[k] = v
		}
	}
}
