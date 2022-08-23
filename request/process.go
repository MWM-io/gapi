package request

import (
	"github.com/mwm-io/gapi/response"
)

// PreProcess is able to execute itself before a server.Handler is executed
type PreProcess interface {
	PreProcess(Handler, *WrappedRequest) (Handler, response.Error)
}

// PreProcessAware is able to give the list of server.PreProcess to execute before executing itself
type PreProcessAware interface {
	PreProcesses() []PreProcess
}

// PreProcessHandler implements the server.PreProcessAware interface
// and is meant to be embedded into the final handler
type PreProcessHandler struct {
	preProcesses []PreProcess
}

// PreProcesses implements the server.PreProcessAware interface
func (p PreProcessHandler) PreProcesses() []PreProcess {
	return p.preProcesses
}

// PPH returns a new PreProcessHandler with the given preProcesses.
func PPH(preProcesses ...PreProcess) PreProcessHandler {
	return PreProcessHandler{
		preProcesses: preProcesses,
	}
}
