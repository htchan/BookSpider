package utils

import (
	"testing"
)

func TestExist(t *testing.T) {
	var testcases = []struct {
		input1   string
		expected bool
	}{
		{"./path_utils_test.go", true},
		{"./not_exist.go", false},
	}
	for _, testcase := range testcases {
		actual := Exists(testcase.input1)
		if actual != testcase.expected {
			t.Fatalf("helper.Exist(\"%v\") result gives\n%v, but not\n%v\n",
				testcase.input1, actual, testcase.expected)
		}
	}
}
