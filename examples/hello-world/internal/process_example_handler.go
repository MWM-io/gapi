package internal

import (
	"github.com/mwm-io/gapi/process"
	"github.com/mwm-io/gapi/request"
	"github.com/mwm-io/gapi/response"
)

// ProcessBody : request body for ProcessHandler
type ProcessBody struct {
	Name string `json:"name"`
}

// ProcessResponse : result for ProcessHandler
type ProcessResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// ProcessHandler /
type ProcessHandler struct {
	request.PreProcessHandler
	request.PostProcessHandler

	body           ProcessBody
	pathParameters struct {
		ID string `path:"id"`
	}
	queryParameters struct {
		Age int `query:"age"`
	}
}

// ProcessHandlerF /
func ProcessHandlerF() request.HandlerFactory {
	return func() request.Handler {
		h := &ProcessHandler{}

		h.PreProcessHandler = request.PreProcessH(
			process.QueryParameters{Parameters: &h.queryParameters},
			process.PathParameters{Parameters: &h.pathParameters},
			process.JsonBody{Body: &h.body},
		)
		h.PostProcessHandler = request.PostProcessH(
			process.SaveRequestID{},
		)
		return h
	}
}

// Serve /
func (h *ProcessHandler) Serve(_ request.WrappedRequest) (interface{}, response.Error) {
	return ProcessResponse{
		ID:   h.pathParameters.ID,
		Name: h.body.Name,
		Age:  h.queryParameters.Age,
	}, nil
}
