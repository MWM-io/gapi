package log

import (
	"os"
	"strings"
)

func ExampleLogger_Log() {
	logger := NewDefaultLogger(NewWriter(
		EntryMarshalerFunc(func(entry Entry) []byte {
			return []byte(strings.Join(
				[]string{
					entry.Message,
					entry.Severity.String(),
				},
				"\n",
			))
		}),
		os.Stdout,
	))

	logger.Log("log message 1")
	logger.Error("log message 2")

	// Output:
	// log message 1
	// default
	// log message 2
	// error
}
