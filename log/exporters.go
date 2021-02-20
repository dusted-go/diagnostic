package log

import "fmt"

// Exporter emits log messages to an output source.
type Exporter interface {
	Export(string)
}

// StdoutExporter emits log events to stdout.
type StdoutExporter struct{}

// Export writes the output directly to stdout.
func (e *StdoutExporter) Export(output string) {
	fmt.Println(output)
}
