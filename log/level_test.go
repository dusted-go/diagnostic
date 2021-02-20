package log

import (
	"testing"
)

func Test_ParseLevel_ReturnsCorrectLevel(t *testing.T) {

	type testCase struct {
		Value    string
		Expected Level
	}

	testCases := []testCase{
		{"0", Default},
		{"100", Debug},
		{"200", Info},
		{"300", Notice},
		{"400", Warning},
		{"500", Error},
		{"600", Critical},
		{"700", Alert},
		{"800", Emergency},
		{"Default", Default},
		{"Debug", Debug},
		{"Info", Info},
		{"Notice", Notice},
		{"Warning", Warning},
		{"Error", Error},
		{"Critical", Critical},
		{"Alert", Alert},
		{"Emergency", Emergency},
		{"default", Default},
		{"debug", Debug},
		{"info", Info},
		{"notice", Notice},
		{"warning", Warning},
		{"error", Error},
		{"critical", Critical},
		{"alert", Alert},
		{"emergency", Emergency},
		{"DEFAULT", Default},
		{"DEBUG", Debug},
		{"INFO", Info},
		{"NOTICE", Notice},
		{"WARNING", Warning},
		{"ERROR", Error},
		{"CRITICAL", Critical},
		{"ALERT", Alert},
		{"EMERGENCY", Emergency},
		{" default ", Default},
		{" debug ", Debug},
		{" info ", Info},
		{" notice ", Notice},
		{" warning ", Warning},
		{" error ", Error},
		{" critical ", Critical},
		{" alert ", Alert},
		{" emergency ", Emergency},
	}

	for _, testCase := range testCases {
		if actual := ParseLevel(testCase.Value); actual != testCase.Expected {
			t.Errorf("\nExpected:\n%s,\nActual:\n%s", testCase.Expected, actual)
		}
	}
}
