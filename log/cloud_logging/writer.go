package cloud_logging

import (
	"github.com/mwm-io/gapi/log"

	"cloud.google.com/go/logging"
	logpb "google.golang.org/genproto/googleapis/logging/v2"
)

// Writer implements the logger.EntryWriter interface
type Writer struct {
	logger *logging.Logger
}

// NewWriter return a new cloud logging logger.EntryWriter.
// See cloud.google.com/go/logging.NewClient on how to create a new logging.Logger.
func NewWriter(log *logging.Logger) *Writer {
	return &Writer{logger: log}
}

// WriteEntry implements the logger.EntryWriter interface
func (w *Writer) WriteEntry(entry log.Entry) {
	loggingEntry := logging.Entry{
		Timestamp: entry.Timestamp,
		Severity:  mapSeverity(entry.Severity),
		Payload:   nil,
		Labels:    entry.Labels,
	}

	if entry.Fields == nil || len(entry.Fields) == 0 {
		loggingEntry.Payload = entry.Message
	} else {
		loggingEntry.Payload = struct {
			Message string                 `json:"message"`
			Fields  map[string]interface{} `json:"fields"`
		}{
			Message: entry.Message,
			Fields:  entry.Fields,
		}
	}

	if entry.TraceID != "" || entry.SpanID != "" {
		loggingEntry.Trace = entry.TraceID
		loggingEntry.SpanID = entry.SpanID
		loggingEntry.TraceSampled = entry.IsTraceSampled
	}

	if entry.StackTrace != nil {
		stackInfo := entry.GetLastStackInfo()

		loggingEntry.SourceLocation = &logpb.LogEntrySourceLocation{
			File:     stackInfo.File,
			Function: stackInfo.Function,
			Line:     int64(stackInfo.Line),
		}
	}

	w.logger.Log(loggingEntry)
}

func mapSeverity(severity log.Severity) logging.Severity {
	var mapping = map[log.Severity]logging.Severity{
		log.DefaultSeverity:   logging.Default,
		log.EmergencySeverity: logging.Emergency,
		log.AlertSeverity:     logging.Alert,
		log.CriticalSeverity:  logging.Critical,
		log.ErrorSeverity:     logging.Error,
		log.WarnSeverity:      logging.Warning,
		log.InfoSeverity:      logging.Info,
		log.DebugSeverity:     logging.Debug,
	}

	loggingSeverity, ok := mapping[severity]
	if !ok {
		return logging.Default
	}

	return loggingSeverity
}
