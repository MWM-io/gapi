package log

import (
	"encoding/json"
	"io"
)

// EntryWriter is able to write an entry.
type EntryWriter interface {
	WriteEntry(Entry)
}

// EntryMarshaler is able to serialize an entry to a []byte.
type EntryMarshaler interface {
	Marshal(Entry) []byte
}

// The EntryMarshalerFunc type is an adapter to allow the use of ordinary functions as EntryMarshaler.
// If f is a function with the appropriate signature, EntryMarshalerFunc(f) is a EntryMarshaler that calls f.
type EntryMarshalerFunc func(entry Entry) []byte

// Marshal implements the EntryMarshaler interface.
func (s EntryMarshalerFunc) Marshal(entry Entry) []byte {
	return s(entry)
}

// JSONEntryMarshaler serialize an entry to JSON.
var JSONEntryMarshaler = EntryMarshalerFunc(func(entry Entry) []byte {
	jsEntry := struct {
		Entry
		StackTrace []StackInfo
	}{
		Entry:      entry,
		StackTrace: entry.GetStackInfo(),
	}

	js, _ := json.Marshal(jsEntry)

	return js
})

// Writer is able to write entries, transforming them given an EntryMarshaler and writing them using an io.Writer.
type Writer struct {
	marshaler EntryMarshaler
	writer    io.Writer
}

// NewWriter returns a new writer.
func NewWriter(marshaler EntryMarshaler, writer io.Writer) *Writer {
	return &Writer{
		marshaler: marshaler,
		writer:    writer,
	}
}

// WriteEntry implements the EntryWriter interface.
func (s *Writer) WriteEntry(entry Entry) {
	_, _ = s.writer.Write(s.marshaler.Marshal(entry))
	_, _ = s.writer.Write([]byte("\n"))
}

// MultiWriter is an EntryWriter that write to all its children.
type MultiWriter struct {
	writers []EntryWriter
}

// NewMultiWriter returns a new MultiWriter.
func NewMultiWriter(writers ...EntryWriter) *MultiWriter {
	return &MultiWriter{writers: writers}
}

// WriteEntry implements the EntryWriter interface.
func (m *MultiWriter) WriteEntry(entry Entry) {
	for _, writer := range m.writers {
		writer.WriteEntry(entry)
	}
}

// FilterWriter is an EntryWriter that will write only logs with the given severity or higher.
type FilterWriter struct {
	severity Severity
	writer   EntryWriter
}

// NewFilterWriter returns a new FilterWriter.
func NewFilterWriter(severity Severity, writer EntryWriter) *FilterWriter {
	return &FilterWriter{
		severity: severity,
		writer:   writer,
	}
}

// WriteEntry implements the EntryWriter interface
func (s *FilterWriter) WriteEntry(entry Entry) {
	if entry.Severity > s.severity {
		return
	}

	s.writer.WriteEntry(entry)
}
