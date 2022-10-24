// Package stacktrace provides you a StackTrace type
// that contains both the list of the runtime.Frame of the stacktrace,
// and the printed format fo debug.Stack(),
// so you can choose which format you want to use.
package stacktrace

import (
	"runtime"
	"runtime/debug"
	"strings"
)

// StackTrace wrap two information:
// - the stacktrace as a string from debug.Stack()
// - the list of the frames coming from runtime.Callers()
type StackTrace struct {
	stack  string
	frames []runtime.Frame
}

// New builds a new StackTrace.
func New() StackTrace {
	pc := make([]uintptr, 32)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frameSlice := make([]runtime.Frame, 0, n)

	for {
		frame, more := frames.Next()

		frameSlice = append(frameSlice, frame)

		if !more {
			break
		}
	}

	return StackTrace{
		stack:  string(debug.Stack()),
		frames: frameSlice,
	}
}

// String returns the stacktrace as a string,
// with the same format as debug.Stack()
func (st StackTrace) String() string {
	return st.stack
}

// Frames return the list of runtime.Frame,
// built using runtime.Callers().
func (st StackTrace) Frames() []runtime.Frame {
	return st.frames
}

// LastFrame will return the last frame in the stacktrace,
// ignoring go source files if ignoreGo is set to true,
// and ignoring files containing the strings provided in ignoredFiles.
// It can happen that all frames are ignored, so you need to check if LastFrame is returning true
func (st StackTrace) LastFrame(ignoreGo bool, ignoredFiles ...string) (runtime.Frame, bool) {
frameLoop:
	for _, frame := range st.frames {
		if ignoreGo {
			if len(strings.Split(frame.Function, "/")) == 1 {
				continue
			}
		}

		for _, ignoredFile := range ignoredFiles {
			if strings.Contains(frame.File, ignoredFile) {
				continue frameLoop
			}
		}

		return frame, true
	}

	return runtime.Frame{}, false
}

// Last return the LastFrame outside the vendors and the go source files.
func (st StackTrace) Last() (runtime.Frame, bool) {
	return st.LastFrame(true, "/vendor/")
}
