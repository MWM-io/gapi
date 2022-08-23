package request

import (
	"github.com/mwm-io/gapi/error"
)

// PreProcess is able to execute itself before a server.Handler is executed
type PreProcess interface {
	PreProcess(Handler, *WrappedRequest) (Handler, error.Error)
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

// PreProcessH returns a new PreProcessHandler with the given preProcesses.
func PreProcessH(preProcesses ...PreProcess) PreProcessHandler {
	return PreProcessHandler{
		preProcesses: preProcesses,
	}
}
