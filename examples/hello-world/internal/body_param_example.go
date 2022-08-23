package internal

import (
	"github.com/mwm-io/gapi/process"
	"github.com/mwm-io/gapi/request"
	"github.com/mwm-io/gapi/response"
)

// ParseParamsBody : request body for ParseParamsHandler
type ParseParamsBody struct {
	Name string `json:"name"`
}

// ParseParamsResponse : result for ParseParamsHandler
type ParseParamsResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// ParseParamsHandler /
type ParseParamsHandler struct {
	request.PreProcessHandler

	body           ParseParamsBody
	pathParameters struct {
		ID string `path:"id"`
	}
	queryParameters struct {
		Age int `query:"age"`
	}
}

// ParseParamsHandlerF /
func ParseParamsHandlerF() request.HandlerFactory {
	return func() request.Handler {
		h := &ParseParamsHandler{}

		h.PreProcessHandler = request.PPH(
			process.QueryParameters{Parameters: &h.queryParameters},
			process.PathParameters{Parameters: &h.pathParameters},
			process.JsonBody{Body: &h.body},
		)
		return h
	}
}

// Serve /
func (h *ParseParamsHandler) Serve(_ request.WrappedRequest) (interface{}, response.Error) {
	return ParseParamsResponse{
		ID:   h.pathParameters.ID,
		Name: h.body.Name,
		Age:  h.queryParameters.Age,
	}, nil
}
