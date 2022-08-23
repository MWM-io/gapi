package request

import (
	"github.com/mwm-io/gapi/error"
)

// PostProcess is able to execute itself before a request.Handler is executed
type PostProcess interface {
	PostProcess(Handler, *WrappedRequest) (Handler, error.Error)
}

// PostProcessAware is able to give the list of request.PostProcess to execute before executing itself
type PostProcessAware interface {
	PostProcesses() []PostProcess
}

// PostProcessHandler implements the request.PostProcessAware interface
// and is meant to be embedded into the final handler
type PostProcessHandler struct {
	preProcesses []PostProcess
}

// PostProcesses implements the request.PostProcessAware interface
func (p PostProcessHandler) PostProcesses() []PostProcess {
	return p.preProcesses
}

// PostProcessH returns a new PostProcessHandler with the given preProcesses.
func PostProcessH(preProcesses ...PostProcess) PostProcessHandler {
	return PostProcessHandler{
		preProcesses: preProcesses,
	}
}
