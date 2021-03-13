package log

import "github.com/dusted-go/diagnostic/trace"

// New creates a new default log event.
func New(filter Filter, formatter Formatter, exporter Exporter, minLevel Level) Event {
	if filter == nil {
		filter = &NoFilter{}
	}

	if formatter == nil {
		formatter = &Console{}
	}

	if exporter == nil {
		exporter = &StdoutExporter{}
	}

	return event{
		filter:         filter,
		formatter:      formatter,
		exporter:       exporter,
		minLevel:       minLevel,
		level:          Debug,
		hasHTTPRequest: false,
	}
}

// NewWithTrace creates a new default log event with initialised trace IDs.
func NewWithTrace(filter Filter, formatter Formatter, exporter Exporter, minLevel Level) Event {
	if filter == nil {
		filter = &NoFilter{}
	}

	if formatter == nil {
		formatter = &Console{}
	}

	if exporter == nil {
		exporter = &StdoutExporter{}
	}

	traceID, spanID := trace.DefaultGenerator.NewTraceIDs()
	return event{
		filter:         filter,
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
	DefaultEvent = New(&NoFilter{}, &Console{}, &StdoutExporter{}, Debug)
)
