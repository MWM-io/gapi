package log

import (
	"fmt"
)

// Severity is the severity level of a message.
// Nothing fancy here, just based on syslog severity levels.
type Severity int

const (
	// DefaultSeverity is the default severity. It means no severity has been set.
	DefaultSeverity Severity = iota
	// EmergencySeverity level: highest level of severity.
	EmergencySeverity
	// AlertSeverity level.Should be corrected ASAP.
	AlertSeverity
	// CriticalSeverity level. Indicates a failure in a primary system.
	CriticalSeverity
	// ErrorSeverity level. Used for errors that should definitely be noted.
	ErrorSeverity
	// WarnSeverity level. Non-critical entries that deserve eyes.
	WarnSeverity
	// InfoSeverity level. Describe normal behavior of the application.
	InfoSeverity
	// DebugSeverity should be used only for debugging purposes. Can contain more verbose information.
	DebugSeverity
)

// String implements the fmt.Stringer interface.
func (s Severity) String() string {
	switch s {
	case DefaultSeverity:
		return "default"
	case EmergencySeverity:
		return "emergency"
	case AlertSeverity:
		return "alert"
	case CriticalSeverity:
		return "critical"
	case ErrorSeverity:
		return "error"
	case WarnSeverity:
		return "warning"
	case InfoSeverity:
		return "info"
	case DebugSeverity:
		return "debug"
	}

	return "unknown"
}

// MarshalJSON implements the json.Marshaler interface.
func (s Severity) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, s.String())), nil
}