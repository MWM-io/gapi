package stacktrace

import (
	"runtime"
	"runtime/debug"
	"strings"
)

type StackTrace struct {
	stack  string
	frames []runtime.Frame
}

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

func (st StackTrace) String() string {
	return st.stack
}

func (st StackTrace) Frames() []runtime.Frame {
	return st.frames
}

func (st StackTrace) LastFrame(ignoredFiles ...string) (runtime.Frame, bool) {
frameLoop:
	for _, frame := range st.frames {
		for _, ignoredFile := range ignoredFiles {
			if strings.Contains(frame.File, ignoredFile) {
				continue frameLoop
			}
		}

		return frame, true
	}

	return runtime.Frame{}, false
}

// Last return the LastFrame outside the gapi package.
func (st StackTrace) Last() (runtime.Frame, bool) {
	return st.LastFrame("mwm-io/gapi")
}
