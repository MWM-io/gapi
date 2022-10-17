package openapi

import (
	"github.com/swaggest/openapi-go/openapi3"
)

type config struct {
	summary     string
	description string
	examples    map[string]openapi3.Example
	mimeType    string
	statusCode  int
	headers     map[string]string
}

func (c *config) applyOptions(options []BuilderOption) {
	for _, option := range options {
		option.apply(c)
	}
}

// BuilderOption can be used to override default value or add precision to documented element
type BuilderOption interface {
	apply(c *config)
}

type withDescription struct {
	description string
}

// WithDescription adds a description to documented element
func WithDescription(description string) BuilderOption {
	return withDescription{
		description: description,
	}
}

func (o withDescription) apply(c *config) {
	c.description = o.description
}

type withExample struct {
	example openapi3.Example
	key     string
}

func (o withExample) apply(c *config) {
	if c.examples == nil {
		c.examples = make(map[string]openapi3.Example)
	}

	c.examples[o.key] = o.example
}

// WithExample add an example to documented element
func WithExample(exampleName string, value interface{}, options ...BuilderOption) BuilderOption {
	var c config
	c.applyOptions(options)

	var description *string
	if c.description != "" {
		description = &c.description
	}

	return withExample{
		key: exampleName,
		example: openapi3.Example{
			Summary:     &exampleName,
			Description: description,
			Value:       &value,
		},
	}
}

type withMimeType struct {
	valueType string
}

func (o withMimeType) apply(c *config) {
	c.mimeType = o.valueType
}

// WithMimeType add precision about type/format (or override default value) to documented element
func WithMimeType(valueType string) BuilderOption {
	return withMimeType{
		valueType: valueType,
	}
}

type withStatusCode struct {
	statusCode int
}

func (o withStatusCode) apply(c *config) {
	c.statusCode = o.statusCode
}

// WithStatusCode override default status code to documented element
func WithStatusCode(statusCode int) BuilderOption {
	return withStatusCode{
		statusCode: statusCode,
	}
}

// Redirect carries redirect information.
type Redirect interface {
	StatusCode() int
	Location() string
}

type withRedirect struct {
	redirect Redirect
}

// WithRedirect add an additional description to documented element
func WithRedirect(redirect Redirect) BuilderOption {
	return withRedirect{
		redirect: redirect,
	}
}

func (o withRedirect) apply(c *config) {
	c.statusCode = o.redirect.StatusCode()

	headers := withHeaders{
		headers: map[string]string{
			"Location": o.redirect.Location(),
		},
	}

	headers.apply(c)
}

type withHeaders struct {
	headers map[string]string
}

// WithHeaders add headers
func WithHeaders(headers map[string]string) BuilderOption {
	return withHeaders{
		headers: headers,
	}
}

func (o withHeaders) apply(c *config) {
	if c.headers == nil {
		c.headers = o.headers
		return
	}

	for k, v := range c.headers {
		c.headers[k] = v
	}
}
