package log

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dusted-go/diagnostic/trace"
)

// Event allows to create and write an event to an output source.
type Event interface {
	SetFormatter(Formatter) Event
	SetExporter(Exporter) Event
	SetMinLogLevel(Level) Event
	SetServiceName(string) Event
	SetServiceVersion(string) Event
	SetHTTPRequest(*http.Request) Event
	SetError(error) Event
	SetData(interface{}) Event
	SetTraceID(trace.ID) Event
	SetSpanID(trace.SpanID) Event
	AddLabel(string, string) Event

	Debug() Event
	Info() Event
	Notice() Event
	Warning() Event
	Error() Event
	Critical() Event
	Alert() Event
	Emergency() Event

	Msg(string)
	Fmt(string, ...interface{})
}

// --------------------------------
// Event implementation
// --------------------------------

type serviceContext struct {
	Name    string
	Version string
}

type httpRequest struct {
	RequestMethod string `json:"requestMethod"`
	RequestURL    string `json:"requestUrl"`
	RequestSize   string `json:"requestSize"`
	UserAgent     string `json:"userAgent"`
	RemoteIP      string `json:"remoteIp"`
	ServerIP      string `json:"serverIp"`
	Referer       string `json:"referer"`
	Protocol      string `json:"protocol"`
}

type event struct {
	filter         Filter
	formatter      Formatter
	exporter       Exporter
	minLevel       Level
	level          Level
	serviceName    string
	serviceVersion string
	httpRequest    httpRequest
	hasHTTPRequest bool
	err            error
	data           interface{}
	labels         map[string]string
	traceID        trace.ID
	spanID         trace.SpanID
	message        string
}

func (e event) SetFilter(filter Filter) Event {
	if filter == nil {
		filter = &NoFilter{}
	}
	e.filter = filter
	return e
}

func (e event) SetFormatter(formatter Formatter) Event {
	if formatter == nil {
		formatter = &Console{}
	}
	e.formatter = formatter
	return e
}

func (e event) SetExporter(exporter Exporter) Event {
	if exporter == nil {
		exporter = &StdoutExporter{}
	}
	e.exporter = exporter
	return e
}

func (e event) SetMinLogLevel(minLevel Level) Event {
	e.minLevel = minLevel
	return e
}

func (e event) SetServiceName(name string) Event {
	e.serviceName = name
	return e
}

func (e event) SetServiceVersion(version string) Event {
	e.serviceVersion = version
	return e
}

func (e event) SetHTTPRequest(req *http.Request) Event {
	if req == nil {
		return e
	}

	reqURL := req.Host + req.RequestURI
	if req.URL != nil {
		reqURL = req.URL.String()
	}

	e.httpRequest = httpRequest{
		RequestMethod: req.Method,
		RequestURL:    reqURL,
		RequestSize:   strconv.FormatInt(req.ContentLength, 10),
		UserAgent:     req.UserAgent(),
		RemoteIP:      req.RemoteAddr,
		Referer:       req.Referer(),
		Protocol:      req.Proto,
		// ToDo: ServerIP
	}
	e.hasHTTPRequest = true
	return e
}

func (e event) SetError(err error) Event {
	e.err = err
	return e
}

func (e event) SetData(data interface{}) Event {
	e.data = data
	return e
}

func (e event) SetTraceID(traceID trace.ID) Event {
	e.traceID = traceID
	return e
}

func (e event) SetSpanID(spanID trace.SpanID) Event {
	e.spanID = spanID
	return e
}

func (e event) AddLabel(key, value string) Event {
	if e.labels == nil {
		e.labels = make(map[string]string)
	}
	e.labels[key] = value
	return e
}

func (e event) setLevel(lvl Level) Event {
	e.level = lvl
	return e
}

// Debug sets the log level as debug.
func (e event) Debug() Event {
	return e.setLevel(Debug)
}

// Info sets the log level as info.
func (e event) Info() Event {
	return e.setLevel(Info)
}

// Notice sets the log level as notice.
func (e event) Notice() Event {
	return e.setLevel(Notice)
}

// Warning sets the log level as warning.
func (e event) Warning() Event {
	return e.setLevel(Warning)
}

// Error sets the log level as error.
func (e event) Error() Event {
	return e.setLevel(Error)
}

// Critical sets the log level as critical.
func (e event) Critical() Event {
	return e.setLevel(Critical)
}

// Alert sets the log level as alert.
func (e event) Alert() Event {
	return e.setLevel(Alert)
}

// Emergency sets the log level as emergency.
func (e event) Emergency() Event {
	return e.setLevel(Emergency)
}

// Msg emits a log event message.
func (e event) Msg(message string) {
	if e.level >= e.minLevel && e.filter.CanWrite(e.message) {
		e.message = message
		e.exporter.Export(e.formatter.Format(e))
	}
}

// Fmt emits a formatted log event message.
func (e event) Fmt(format string, args ...interface{}) {
	if e.level >= e.minLevel {
		e.message = fmt.Sprintf(format, args...)
		if e.filter.CanWrite(e.message) {
			e.exporter.Export(e.formatter.Format(e))
		}
	}
}
