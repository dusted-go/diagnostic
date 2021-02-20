package log

import (
	"strconv"
	"strings"
)

// Level denotes the importance of a log event.
type Level int

const (
	// Default means that the log event has no assigned log level.
	Default Level = iota * 100
	// Debug or trace information.
	Debug
	// Info events are routine information, such as ongoing status or performance.
	Info
	// Notice are normal but significant events, such as start up, shut down, or a configuration change.
	Notice
	// Warning events might cause problems.
	Warning
	// Error events are likely to cause problems.
	Error
	// Critical events cause more severe problems or outages.
	Critical
	// Alert events mean that a person must take an action immediately.
	Alert
	// Emergency means that one or more systems are unusable.
	Emergency
)

func (lvl Level) String() string {
	switch lvl {
	case Debug:
		return "DEBUG"
	case Info:
		return "INFO"
	case Notice:
		return "NOTICE"
	case Warning:
		return "WARNING"
	case Error:
		return "ERROR"
	case Critical:
		return "CRITICAL"
	case Alert:
		return "ALERT"
	case Emergency:
		return "EMERGENCY"
	default:
		return "DEFAULT"
	}
}

// Short returns a 3 letter acronym for the log level.
func (lvl Level) Short() string {
	switch lvl {
	case Debug:
		return "DBG"
	case Info:
		return "INF"
	case Notice:
		return "NTC"
	case Warning:
		return "WRN"
	case Error:
		return "ERR"
	case Critical:
		return "CRT"
	case Alert:
		return "ALR"
	case Emergency:
		return "EMG"
	default:
		return "DFT"
	}
}

// ParseLevel parses a string value into a log level.
func ParseLevel(value string) Level {
	if len(value) > 0 {
		v := strings.ToLower(strings.Trim(value, " "))

		if v == "default" {
			return Default
		}
		if v == "debug" {
			return Debug
		}
		if v == "info" {
			return Info
		}
		if v == "notice" {
			return Notice
		}
		if v == "warning" {
			return Warning
		}
		if v == "error" {
			return Error
		}
		if v == "critical" {
			return Critical
		}
		if v == "alert" {
			return Alert
		}
		if v == "emergency" {
			return Emergency
		}

		if i, err := strconv.Atoi(value); err == nil {
			return Level(i)
		}
	}
	return Default
}
