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

type CoreOption func(*CoreConfig)

func WithTracing(tracing Tracing) CoreOption {
	return func(config *CoreConfig) {
		config.Tracing = tracing
	}
}

func WithLog(log Log) CoreOption {
	return func(config *CoreConfig) {
		config.Log = log
	}
}

func WithLogResponseWriter(log Log) CoreOption {
	return func(config *CoreConfig) {
		config.LogResponseWriter = log
	}
}

func WithResponseType(response interface{}) CoreOption {
	return func(config *CoreConfig) {
		config.ResponseWriter.Response = response
	}
}

func WithResponseMarshaler(contentType string, marshaler Marshaler) CoreOption {
	return func(config *CoreConfig) {
		config.ResponseWriter.Marshalers[contentType] = marshaler
	}
}

func WithBodyUnmarshaler(contentType string, unmarshaler Unmarshaler) CoreOption {
	return func(config *CoreConfig) {
		config.BodyUnmarshaler.Unmarshalers[contentType] = unmarshaler
	}
}

func WithDefaultContentType(contentType string) CoreOption {
	return func(config *CoreConfig) {
		config.ResponseWriter.DefaultContentType = contentType
		config.BodyUnmarshaler.DefaultContentType = contentType
	}
}

func WithPathParameters(params interface{}) CoreOption {
	return func(config *CoreConfig) {
		config.PathParameters.Parameters = params
	}
}

func WithQueryParameters(params interface{}) CoreOption {
	return func(config *CoreConfig) {
		config.QueryParameters.Parameters = params
	}
}

func WithBody(body interface{}) CoreOption {
	return func(config *CoreConfig) {
		config.BodyUnmarshaler.Body = body
	}
}

func WithSkipBodyValidation(skip bool) CoreOption {
	return func(config *CoreConfig) {
		config.BodyUnmarshaler.SkipValidation = skip
	}
}
