package log

import (
	"testing"
)

type testCase struct {
	Level    Level
	Msg      string
	Err      error
	Data     interface{}
	Expected string
}

type Address struct {
	HouseNumber int
	Street      string
	Postcode    string
}

type Person struct {
	FirstName string
	LastName  string
	Pets      []string
	Age       int
	Address   Address
}

func Test_Stackdriver_WithoutAnyExtraSettings_FormatsCorrectly(t *testing.T) {
	stackdriver := Stackdriver{}

	testCases := []testCase{
		{Debug, "debug", nil, nil, "{\"severity\":\"DEBUG\",\"message\":\"debug\"}"},
		{Info, "info", nil, nil, "{\"severity\":\"INFO\",\"message\":\"info\"}"},
		{Notice, "notice", nil, nil, "{\"severity\":\"NOTICE\",\"message\":\"notice\"}"},
		{Warning, "warning", nil, nil, "{\"severity\":\"WARNING\",\"message\":\"warning\"}"},
		{Error, "error", nil, nil, "{\"severity\":\"ERROR\",\"message\":\"error\"}"},
		{Critical, "critical", nil, nil, "{\"severity\":\"CRITICAL\",\"message\":\"critical\"}"},
		{Alert, "alert", nil, nil, "{\"severity\":\"ALERT\",\"message\":\"alert\"}"},
		{Emergency, "emergency", nil, nil, "{\"severity\":\"EMERGENCY\",\"message\":\"emergency\"}"},

		{Alert, "foo bar", nil, Address{HouseNumber: 1, Street: "Foo", Postcode: "Bar"}, "{\"severity\":\"ALERT\",\"message\":\"foo bar\",\"data\":{\"HouseNumber\":1,\"Street\":\"Foo\",\"Postcode\":\"Bar\"}}"},

		{Info, "yada", nil, Person{FirstName: "Sue", LastName: "Doe", Age: 45, Address: Address{HouseNumber: 3, Street: "x", Postcode: "Y"}}, "{\"severity\":\"INFO\",\"message\":\"yada\",\"data\":{\"FirstName\":\"Sue\",\"LastName\":\"Doe\",\"Pets\":null,\"Age\":45,\"Address\":{\"HouseNumber\":3,\"Street\":\"x\",\"Postcode\":\"Y\"}}}"},
	}

	for _, testCase := range testCases {
		if actual := stackdriver.Format(
			event{
				level:   testCase.Level,
				message: testCase.Msg,
				err:     testCase.Err,
				data:    testCase.Data,
			}); actual != testCase.Expected {
			t.Errorf("\nExpected:\n%s,\nActual:\n%s", testCase.Expected, actual)
		}
	}
}

func Test_Stackdriver_WithServiceContext_FormatsCorrectly(t *testing.T) {
	stackdriver := Stackdriver{}

	data := Person{
		FirstName: "Sue",
		LastName:  "Doe",
		Age:       45,
		Address: Address{
			HouseNumber: 3,
			Street:      "x",
			Postcode:    "Y"}}

	expected := "{\"severity\":\"INFO\",\"message\":\"this is a stupid message\",\"serviceContext.service\":\"foo-bar\",\"serviceContext.version\":\"v1.0.0\",\"data\":{\"FirstName\":\"Sue\",\"LastName\":\"Doe\",\"Pets\":null,\"Age\":45,\"Address\":{\"HouseNumber\":3,\"Street\":\"x\",\"Postcode\":\"Y\"}}}"

	if actual := stackdriver.Format(event{
		serviceName:    "foo-bar",
		serviceVersion: "v1.0.0",
		level:          Info,
		message:        "this is a stupid message",
		err:            nil,
		data:           data,
	}); actual != expected {
		t.Errorf("\nExpected:\n%s,\nActual:\n%s", expected, actual)
	}
}

func Test_Stackdriver_WithLabels_FormatsCorrectly(t *testing.T) {
	stackdriver := Stackdriver{}

	data := Person{
		FirstName: "Sue",
		LastName:  "Doe",
		Age:       45,
		Address: Address{
			HouseNumber: 3,
			Street:      "x",
			Postcode:    "Y"}}

	expected := "{\"severity\":\"INFO\",\"message\":\"this is a stupid message\",\"logging.googleapis.com/labels\":{\"B\":\"b\",\"a\":\"A\"},\"data\":{\"FirstName\":\"Sue\",\"LastName\":\"Doe\",\"Pets\":null,\"Age\":45,\"Address\":{\"HouseNumber\":3,\"Street\":\"x\",\"Postcode\":\"Y\"}}}"

	if actual := stackdriver.Format(event{
		level:   Info,
		message: "this is a stupid message",
		err:     nil,
		data:    data,
		labels: map[string]string{
			"a": "A",
			"B": "b",
		},
	}); actual != expected {
		t.Errorf("\nExpected:\n%s,\nActual:\n%s", expected, actual)
	}
}

func Test_Stackdriver_WithLabelsAndServiceContext_FormatsCorrectly(t *testing.T) {
	stackdriver := Stackdriver{}

	data := Person{
		FirstName: "Sue",
		LastName:  "Doe",
		Age:       45,
		Address: Address{
			HouseNumber: 3,
			Street:      "x",
			Postcode:    "Y"}}

	expected := "{\"severity\":\"INFO\",\"message\":\"this is a stupid message\",\"serviceContext.service\":\"foo-bar\",\"serviceContext.version\":\"v1.0.0\",\"logging.googleapis.com/labels\":{\"B\":\"b\",\"a\":\"A\"},\"data\":{\"FirstName\":\"Sue\",\"LastName\":\"Doe\",\"Pets\":null,\"Age\":45,\"Address\":{\"HouseNumber\":3,\"Street\":\"x\",\"Postcode\":\"Y\"}}}"

	if actual := stackdriver.Format(event{
		serviceName:    "foo-bar",
		serviceVersion: "v1.0.0",
		level:          Info,
		message:        "this is a stupid message",
		err:            nil,
		data:           data,
		labels: map[string]string{
			"a": "A",
			"B": "b",
		},
	}); actual != expected {
		t.Errorf("\nExpected:\n%s,\nActual:\n%s", expected, actual)
	}
}

func Test_Stackdriver_WithLabelsAndServiceContextAndHttpRequest_FormatsCorrectly(t *testing.T) {
	stackdriver := Stackdriver{}

	data := Person{
		FirstName: "Sue",
		LastName:  "Doe",
		Age:       45,
		Address: Address{
			HouseNumber: 3,
			Street:      "x",
			Postcode:    "Y"}}

	expected := "{\"severity\":\"INFO\",\"message\":\"this is a stupid message\",\"serviceContext.service\":\"foo-bar\",\"serviceContext.version\":\"v1.0.0\",\"logging.googleapis.com/labels\":{\"B\":\"b\",\"a\":\"A\"},\"httpRequest\":{\"requestMethod\":\"GET\",\"requestUrl\":\"http://example.org/\",\"requestSize\":\"132\",\"userAgent\":\"abc\",\"remoteIp\":\"127.0.0.1\",\"serverIp\":\"\",\"referer\":\"google.com\",\"protocol\":\"HTTP/1.1\"},\"data\":{\"FirstName\":\"Sue\",\"LastName\":\"Doe\",\"Pets\":null,\"Age\":45,\"Address\":{\"HouseNumber\":3,\"Street\":\"x\",\"Postcode\":\"Y\"}}}"

	if actual := stackdriver.Format(event{
		serviceName:    "foo-bar",
		serviceVersion: "v1.0.0",
		hasHTTPRequest: true,
		httpRequest: httpRequest{
			RequestMethod: "GET",
			RequestURL:    "http://example.org/",
			RequestSize:   "132",
			UserAgent:     "abc",
			RemoteIP:      "127.0.0.1",
			ServerIP:      "",
			Referer:       "google.com",
			Protocol:      "HTTP/1.1",
		},
		level:   Info,
		message: "this is a stupid message",
		err:     nil,
		data:    data,
		labels: map[string]string{
			"a": "A",
			"B": "b",
		},
	}); actual != expected {
		t.Errorf("\nExpected:\n%s,\nActual:\n%s", expected, actual)
	}
}
