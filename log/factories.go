package log

import "github.com/dusted-go/diagnostics/trace"

// New creates a new default log event.
func New(formatter Formatter, exporter Exporter, minLevel Level) Event {
	return event{
		formatter:      formatter,
		exporter:       exporter,
		minLevel:       minLevel,
		level:          Debug,
		hasHTTPRequest: false,
	}
}

// NewWithTrace creates a new default log event with initialised trace IDs.
func NewWithTrace(formatter Formatter, exporter Exporter, minLevel Level) Event {
	traceID, spanID := trace.DefaultGenerator.NewTraceIDs()
	return event{
		formatter:      formatter,
		exporter:       exporter,
		minLevel:       minLevel,
		level:          Debug,
		hasHTTPRequest: false,
		traceID:        traceID,
		spanID:         spanID,
	}
}

var (
	// DefaultEvent returns a default log event.
	DefaultEvent = New(&Console{}, &StdoutExporter{}, Debug)
)
