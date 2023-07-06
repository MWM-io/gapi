package errors

import (
	"fmt"
	"runtime"
	"strings"
)

// GetCallers return the caller of the function and the call stack
func GetCallers() (callerName, caller string, callStack []string) {
	// Ask runtime.Callers for up to 10 pcs
	pc := make([]uintptr, 10)
	n := runtime.Callers(1, pc)
	if n == 0 {
		// No pcs available. Stop now.
		// This can happen if the first argument to runtime.Callers are large.
		caller = "unknown"
		return
	}

	pc = pc[:n] // pass only valid pcs to runtime.CallersFrames
	frames := runtime.CallersFrames(pc)

	firstFrame := true
	// Loop to get frames.
	// A fixed number of pcs can expand to an indefinite number of Frames.
	for {
		frame, more := frames.Next()
		if !more {
			break
		}

		// Stop call stack when we reach the handler caller
		// every call after this is not relevant
		if !firstFrame && strings.Contains(frame.File, "github.com/mwm-io/gapi") && strings.Contains(frame.File, "/handler/") {
			break
		}

		// Ignore errors package from call trace because all errors was created from this package
		if strings.Contains(frame.File, "github.com/mwm-io/gapi") && strings.Contains(frame.File, "/errors/") {
			continue
		}

		if firstFrame {
			caller = formatFrame(frame)
			callerName = frame.Func.Name()
			firstFrame = false
		} else {
			callStack = append(callStack, formatFrame(frame))
		}
	}

	return
}

func formatFrame(frame runtime.Frame) string {
	file := frame.File
	line := frame.Line
	function := frame.Function

	return fmt.Sprintf("%s:%s -> %s", file, itoa(line, -1), function)
}

// Cheap integer to fixed-width decimal ASCII. Give a negative width to avoid zero-padding.
func itoa(i int, wid int) string {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	return string(b[bp:])
}
