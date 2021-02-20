package trace

import (
	"testing"
)

func Test_SpanID(t *testing.T) {
	headerValue := "2205310701640571284"
	spanID, err := ParseGoogleCloudSpanID(headerValue)

	if err != nil {
		t.Errorf("Test failed due to parsing error:\n%s", err.Error())
		return
	}

	actual := spanID.Decimal()
	expected := uint64(2205310701640571284)

	if err != nil || actual != expected {
		t.Errorf("\nExpected:\n%d,\nActual:\n%d", expected, actual)
	}
}
