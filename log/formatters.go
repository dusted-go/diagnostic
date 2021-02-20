package log

import (
	"encoding/json"
	"fmt"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Formatter can format a log event into a string.
type Formatter interface {
	Format(event) string
}

// --------------------------------
// Console
// --------------------------------

const (
	timeFormat = "2006-01-02 15:04:05.000"

	reset = "\033[0m"

	normal    = 0
	bold      = 2
	extraBold = 4

	black        = 30
	red          = 31
	green        = 32
	yellow       = 33
	blue         = 34
	magenta      = 35
	cyan         = 36
	lightGray    = 37
	darkGray     = 90
	lightRed     = 91
	lightGreen   = 92
	lightYellow  = 93
	lightBlue    = 94
	lightMagenta = 95
	lightCyan    = 96
	white        = 97
)

func logFmt(fontWeight, colorCode int) string {
	return fmt.Sprintf("\033[%s;%sm", strconv.Itoa(fontWeight), strconv.Itoa(colorCode))
}

func logColor(lvl Level) (string, string) {
	switch lvl {
	case Debug:
		return logFmt(normal, darkGray), logFmt(normal, darkGray)
	case Info:
		return logFmt(normal, lightGray), reset
	case Notice:
		return logFmt(normal, lightGreen), reset
	case Warning:
		return logFmt(normal, lightYellow), reset
	case Error:
		return logFmt(normal, lightRed), reset
	case Alert:
		return logFmt(normal, red), reset
	case Critical:
		return logFmt(normal, lightRed), logFmt(normal, lightRed)
	case Emergency:
		return logFmt(normal, red), logFmt(normal, red)
	case Default:
		return logFmt(normal, white), reset
	default:
		return logFmt(normal, white), reset
	}
}

func logLevel(lvl Level) string {
	severityColor, textColor := logColor(lvl)
	return fmt.Sprintf("%s[%s]%s", severityColor, lvl.Short(), textColor)
}

// Console formats an event into a colour formatted human readable text.
type Console struct {
}

// Format formats a log event into the Stackdriver specific JSON schema.
func (f *Console) Format(e event) string {
	errMsg := ""
	if e.err != nil {
		errMsg = fmt.Sprintf("\n\n%s\n\n%s", e.err.Error(), debug.Stack())
	}

	return fmt.Sprintf(
		"%s[%s]%s %s[%s] %s %s%s%s",
		logFmt(normal, blue),
		time.Now().UTC().Format(timeFormat),
		reset,
		logFmt(normal, lightGray),
		e.traceID.String(),
		logLevel(e.level),
		e.message,
		errMsg,
		reset)
}

// --------------------------------
// Stackdriver
// --------------------------------

// Stackdriver formats an event into the Stackdriver specific JSON format.
type Stackdriver struct {
}

func sortKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return keys
}

func escapeJSON(str string) string {
	b, err := json.Marshal(str)
	if err != nil {
		panicOnErr(err, "Failed to escape JSON string.")
	}
	s := string(b)
	return s[1 : len(s)-1]
}

// Format formats a log event into the Stackdriver specific JSON schema.
func (f *Stackdriver) Format(e event) string {

	var str strings.Builder
	str.WriteString("{")
	str.WriteString(fmt.Sprintf("\"severity\":\"%s\"", e.level.String()))

	if e.err == nil {
		str.WriteString(fmt.Sprintf(",\"message\":\"%s\"", escapeJSON(e.message)))
	} else {
		str.WriteString(",\"@type\":\"type.googleapis.com/google.devtools.clouderrorreporting.v1beta1.ReportedErrorEvent\"")
		errMsg := fmt.Sprintf("%+v\n\n%s", e.err.Error(), debug.Stack())
		if len(e.message) > 0 {
			errMsg = e.message + "\n\nError:\n\n" + errMsg
		}
		str.WriteString(fmt.Sprintf(",\"message\":\"%s\"", escapeJSON(errMsg)))
	}

	if e.traceID.IsValid() {
		str.WriteString(fmt.Sprintf(",\"logging.googleapis.com/trace_sampled\":\"true\",\"logging.googleapis.com/trace\":\"%s\"", e.traceID.String()))

		if e.spanID.IsValid() {
			str.WriteString(fmt.Sprintf(",\"logging.googleapis.com/spanId\":\"%d\"", e.spanID.Decimal()))
		}
	}

	if len(e.serviceName) > 0 {
		str.WriteString(fmt.Sprintf(",\"serviceContext.service\":\"%s\"", e.serviceName))
	}

	if len(e.serviceVersion) > 0 {
		str.WriteString(fmt.Sprintf(",\"serviceContext.version\":\"%s\"", e.serviceVersion))
	}

	if e.labels != nil && len(e.labels) > 0 {
		str.WriteString(",\"logging.googleapis.com/labels\":{")

		isFirst := true
		sortedKeys := sortKeys(e.labels)

		for _, key := range sortedKeys {
			value := e.labels[key]
			if !isFirst {
				str.WriteString(",")
			}
			str.WriteString(fmt.Sprintf("\"%s\":\"%s\"", key, value))
			isFirst = false
		}

		str.WriteString("}")
	}

	if e.hasHTTPRequest {
		buffer, err := json.Marshal(e.httpRequest)
		if err == nil {
			req := string(buffer)
			str.WriteString(fmt.Sprintf(",\"httpRequest\":%s", req))
		}
	}

	if e.data != nil {
		var dataStr string
		buffer, err := json.Marshal(e.data)
		if err != nil {
			dataStr = fmt.Sprintf("\"Could not successfully serialize data object into JSON when writing this message.\n\nError: %+v\"", err)
		}
		dataStr = string(buffer)
		str.WriteString(fmt.Sprintf(",\"data\":%s", dataStr))
	}

	str.WriteString("}")
	return str.String()
}
