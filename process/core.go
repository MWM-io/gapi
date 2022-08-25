package process

import (
	"encoding/json"
	"encoding/xml"

	"contrib.go.opencensus.io/exporter/stackdriver/propagation"

	gLog "github.com/mwm-io/gapi/log"
	"github.com/mwm-io/gapi/request"
)

// Core allows you to easily add middlewares to your handler.
func Core(opts ...CoreOption) request.MiddlewareHandler {
	core := CoreConfig{
		Tracing: Tracing{
			Propagation:      &propagation.HTTPFormat{},
			IsPublicEndpoint: false,
		},
		LogResponseWriter: Log{Logger: gLog.GlobalLogger()},
		Log:               Log{Logger: gLog.GlobalLogger()},
		ResponseWriter: request.ResponseWriterMiddleware{
			Marshalers: map[string]request.Marshaler{
				"application/json": request.MarshalerFunc(json.Marshal),
				"application/xml":  request.MarshalerFunc(xml.Marshal),
			},
			DefaultContentType: "application/json",
		},
		PathParameters:  PathParameters{},
		QueryParameters: QueryParameters{},
		JsonBody:        JsonBody{},
	}

	for _, opt := range opts {
		opt(&core)
	}

	return request.MiddlewareH(
		core.Tracing,
		core.LogResponseWriter,
		core.ResponseWriter,
		core.Log,
		core.PathParameters,
		core.QueryParameters,
		core.JsonBody,
	)
}

type CoreConfig struct {
	Tracing           Tracing
	LogResponseWriter Log
	Log               Log
	ResponseWriter    request.ResponseWriterMiddleware
	PathParameters    PathParameters
	QueryParameters   QueryParameters
	JsonBody          JsonBody
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

func WithResponseMarshaler(contentType string, marshaler request.Marshaler) CoreOption {
	return func(config *CoreConfig) {
		config.ResponseWriter.Marshalers[contentType] = marshaler
	}
}

func WithDefaultContentType(contentType string) CoreOption {
	return func(config *CoreConfig) {
		config.ResponseWriter.DefaultContentType = contentType
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
		config.JsonBody.Body = body
	}
}

func WithSkipBodyValidation(skip bool) CoreOption {
	return func(config *CoreConfig) {
		config.JsonBody.SkipValidation = skip
	}
}
