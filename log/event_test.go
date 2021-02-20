package log

import (
	"net/http"
	"testing"
)

func Test_Event_ImmutabilityTest(t *testing.T) {

	httpReq := &http.Request{
		Method:        "GET",
		Host:          "example.org",
		RequestURI:    "/",
		ContentLength: 132,
		Header: map[string][]string{
			"User-Agent": {"abc"},
			"Referer":    {"google.com"},
		},
		RemoteAddr: "127.0.0.1",
		Proto:      "HTTP/1.1",
	}

	data := Person{
		FirstName: "Sue",
		LastName:  "Doe",
		Age:       45,
		Address: Address{
			HouseNumber: 3,
			Street:      "x",
			Postcode:    "Y"}}

	event1 := event{}
	event2 := event1.
		SetHTTPRequest(httpReq).
		SetServiceName("foo-bar").
		SetServiceVersion("v1.0.0").
		AddLabel("a", "A").
		AddLabel("B", "b").
		SetData(data)

	httpReq.Method = "POST"
	data.FirstName = "Linda"

	if event1.data != nil || event1.serviceName != "" || event1.serviceVersion != "" || len(event1.labels) != 0 || event1.hasHTTPRequest {
		t.Error("Log event has been illegally mutated.")
	}

	if event2.(event).data.(Person).FirstName != "Sue" || event2.(event).httpRequest.RequestMethod != "GET" || len(event2.(event).labels) != 2 {
		t.Error("Log event has been illegally mutated.")
	}
}
