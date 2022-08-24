package mattermost

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"time"

	"github.com/mwm-io/gapi/log"
)

// Writer is a log.EntryWriter that will send the log to mattermost on a specific channel.
type Writer struct {
	http      http.Client
	hookURL   string
	channel   string
	onError   func(error)
	formatter func(log.Entry) string
}

// WriterOpt is able to configure a Writer.
type WriterOpt func(writer *Writer)

// OnError set a callback function to execute everytime an error happens when sending a log to mattermost.
func OnError(callback func(error)) WriterOpt {
	return func(writer *Writer) {
		writer.onError = callback
	}
}

// Channel set a specific channel where the logs will be sent.
func Channel(channel string) WriterOpt {
	return func(writer *Writer) {
		writer.channel = channel
	}
}

// Formatter set an entry formatter to format the message sent to mattermost.
// If None
func Formatter(formatter func(log.Entry) string) WriterOpt {
	return func(writer *Writer) {
		writer.formatter = formatter
	}
}

// NewWriter returns a new writer.
// When no channel is given, it will send the message to the default channel (general).
// When no formatter is given, it will take the default formatter.
func NewWriter(http http.Client, hookURL string, opts ...WriterOpt) *Writer {
	w := &Writer{
		http:    http,
		hookURL: hookURL,
	}

	for _, opt := range opts {
		opt(w)
	}

	return w
}

// WriteEntry implements the log.EntryWriter interface.
func (w *Writer) WriteEntry(entry log.Entry) {
	go func() {
		formatter := DefaultFormatter
		if w.formatter != nil {
			formatter = w.formatter
		}
		err := w.sendMessage(formatter(entry))
		if err != nil && w.onError != nil {
			w.onError(err)
		}
	}()
}

func (w *Writer) sendMessage(text string) error {
	messagePayload := struct {
		Channel string `json:"channel"`
		Text    string `json:"text"`
	}{
		Channel: w.channel,
		Text:    text,
	}

	payload, err := json.Marshal(messagePayload)
	if err != nil {
		return err
	}

	resp, err := w.http.Post(w.hookURL, "application/json", bytes.NewReader(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorMessage := fmt.Sprintf("unable to send mattermost message, response status %d", resp.StatusCode)
		respBody, errRead := ioutil.ReadAll(resp.Body)
		if errRead == nil {
			errorMessage = fmt.Sprintf("%s (body: %s)", errorMessage, respBody)
		}

		return errors.New(errorMessage)
	}

	return nil
}

// DefaultFormatter is the default entry formatter used.
func DefaultFormatter(entry log.Entry) string {
	var labelsStr string
	for key, value := range entry.Labels {
		labelsStr = fmt.Sprintf("%s - %s: %s\n", labelsStr, key, value)
	}

	var fieldsStr string
	for key, value := range entry.Fields {
		fieldsStr = fmt.Sprintf("%s - %s: %v\n", fieldsStr, key, value)
	}

	var stackTraceStr string
	frames := runtime.CallersFrames(entry.StackTrace)
	for {
		frame, more := frames.Next()

		stackTraceStr = fmt.Sprintf("%s1. File: %s, Function: %s, Line: %s", stackTraceStr, frame.File, frame.Function, frame.Line)
		if !more {
			break
		}
	}

	return fmt.Sprintf(`
## %s - %s
**%s**

### Labels
%s

### Fields
%s

### Trace
 - TraceID: %s
 - SpanID: %s
 - IsSampled: %v

### StackTrace
%s
`,
		entry.Severity, entry.Timestamp.Format(time.RFC3339),
		entry.Message,
		labelsStr,
		fieldsStr,
		entry.TraceID,
		entry.SpanID,
		entry.IsTraceSampled,
		stackTraceStr,
	)
}
