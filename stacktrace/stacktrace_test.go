package stacktrace

import (
	"fmt"
)

func ExampleLastFrame() {
	a := New()

	// This test file is ignored and we include go source files.
	lastFrame, _ := a.LastFrame(false, "stacktrace_test")
	fmt.Println(lastFrame.Function)

	// we ignore go source files and include this file
	lastFrame, _ = a.LastFrame(true)
	fmt.Println(lastFrame.Function)

	// Ignoring go files here is not mandatory as this file is above in the stack trace.
	lastFrame, _ = a.LastFrame(false)
	fmt.Println(lastFrame.Function)

	// Output:
	// testing.runExample
	// github.com/mwm-io/gapi/stacktrace.ExampleLastFrame
	// github.com/mwm-io/gapi/stacktrace.ExampleLastFrame
}
