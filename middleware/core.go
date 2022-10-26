package middleware

import (
	"encoding/json"
	"encoding/xml"

	"contrib.go.opencensus.io/exporter/stackdriver/propagation"

	gLog "github.com/mwm-io/gapi/log"
	"github.com/mwm-io/gapi/server"
)

// Core allows you to easily add middlewares to your handler.
func Core(opts ...CoreOption) server.MiddlewareHandler {
	core := CoreConfig{
		Tracing: Tracing{
			Propagation:      &propagation.HTTPFormat{},
			IsPublicEndpoint: false,
		},
		LogResponseWriter: Log{Logger: gLog.GlobalLogger()},
		ResponseWriter: ResponseWriterMiddleware{
			Marshalers: map[string]Marshaler{
				"application/json": MarshalerFunc(json.Marshal),
				"application/xml":  MarshalerFunc(xml.Marshal),
			},
			DefaultContentType: "application/json",
		},
		Log:             Log{Logger: gLog.GlobalLogger()},
		PathParameters:  PathParameters{},
		QueryParameters: QueryParameters{},
		BodyUnmarshaler: BodyUnmarshaler{
			Unmarshalers: map[string]Unmarshaler{
				"application/json": UnmarshalerFunc(json.Unmarshal),
				"application/xml":  UnmarshalerFunc(xml.Unmarshal),
			},
			DefaultContentType: "application/json",
		},
	}

	for _, opt := range opts {
		opt(&core)
	}

	return server.MiddlewareH(
		core.Tracing,
		core.LogResponseWriter,
		core.ResponseWriter,
		core.Log,
		core.Recover,
		core.PathParameters,
		core.QueryParameters,
		core.BodyUnmarshaler,
	)
}

// CoreConfig contains all the configuration for the core.
type CoreConfig struct {
	Tracing           Tracing
	LogResponseWriter Log
	ResponseWriter    ResponseWriterMiddleware
	Log               Log
	Recover           Recover
	PathParameters    PathParameters
	QueryParameters   QueryParameters
	BodyUnmarshaler   BodyUnmarshaler
}

// CoreOption is a function that is able to modify the CoreConfig.
type CoreOption func(*CoreConfig)

// WithTracing adds tracing configuration.
func WithTracing(tracing Tracing) CoreOption {
	return func(config *CoreConfig) {
		config.Tracing = tracing
	}
}

// WithLog adds logging configuration.
func WithLog(log Log) CoreOption {
	return func(config *CoreConfig) {
		config.Log = log
	}
}

// WithLogResponseWriter adds logging configuration for the response writer..
func WithLogResponseWriter(log Log) CoreOption {
	return func(config *CoreConfig) {
		config.LogResponseWriter = log
	}
}

// WithResponseType specify the response type. (for the documentation)
func WithResponseType(response interface{}) CoreOption {
	return func(config *CoreConfig) {
		config.ResponseWriter.Response = response
	}
}

// WithResponseMarshaler adds a new response marshaler.
func WithResponseMarshaler(contentType string, marshaler Marshaler) CoreOption {
	return func(config *CoreConfig) {
		config.ResponseWriter.Marshalers[contentType] = marshaler
	}
}

// WithBodyUnmarshaler adds a new body unmarshaler.
func WithBodyUnmarshaler(contentType string, unmarshaler Unmarshaler) CoreOption {
	return func(config *CoreConfig) {
		config.BodyUnmarshaler.Unmarshalers[contentType] = unmarshaler
	}
}

// WithDefaultContentType sets the default content type. (for both the request and the response)
func WithDefaultContentType(contentType string) CoreOption {
	return func(config *CoreConfig) {
		config.ResponseWriter.DefaultContentType = contentType
		config.BodyUnmarshaler.DefaultContentType = contentType
	}
}

// WithForcedContentType sets the content type that will always be returned with the response.
func WithForcedContentType(contentType string) CoreOption {
	return func(config *CoreConfig) {
		config.ResponseWriter.ForcedContentType = contentType
	}
}

// WithPathParameters set the request parameters to populate.
func WithPathParameters(params interface{}) CoreOption {
	return func(config *CoreConfig) {
		config.PathParameters.Parameters = params
	}
}

// WithQueryParameters set the query parameters to populate.
func WithQueryParameters(params interface{}) CoreOption {
	return func(config *CoreConfig) {
		config.QueryParameters.Parameters = params
	}
}

// WithBody set the body to unmarshal the request's body in.
func WithBody(body interface{}) CoreOption {
	return func(config *CoreConfig) {
		config.BodyUnmarshaler.Body = body
	}
}

// WithSkipBodyValidation indicates if we should skip body validation.
func WithSkipBodyValidation(skip bool) CoreOption {
	return func(config *CoreConfig) {
		config.BodyUnmarshaler.SkipValidation = skip
	}
}
