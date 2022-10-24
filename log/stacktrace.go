package log

import (
	"fmt"
	"runtime"
)

// StackTrace represents a go stack-trace.
// It should be able to both return the stacktrace as a string from the debug.Stack() function
// and the list of frames, as well as the last frame to display
// This allows the EntryWriter to choose which format they want.
type StackTrace interface {
	fmt.Stringer
	Frames() []runtime.Frame
	Last() (runtime.Frame, bool)
}
