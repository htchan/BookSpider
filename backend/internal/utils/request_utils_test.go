package utils

import (
	"golang.org/x/text/encoding"
	"testing"
)

func Test_getWeb(t *testing.T) {
	var testcases = []struct {
		input, expected string
	}{
		{"idk", ""},
		{"http://www.google.com/idk", "404"},
		{"https://api.github.com/repos/htchan/Bookspider/downloads", "[]"},
	}
	for _, testcase := range testcases {
		actual := getWeb(testcase.input)
		if actual != testcase.expected {
			t.Fatalf("utils.getWeb(\"%v\") result gives\n\"%v\", but not\n\"%v\"\n",
				testcase.input, actual, testcase.expected)
		}
	}
}
func TestGetWeb(t *testing.T) {
	var testcases = []struct {
		input1    string
		input2    int
		input3    *encoding.Decoder
		input4    int
		expected1 string
		expected2 int
	}{
		{"https://api.github.com/repos/htchan/Bookspider/downloads", 1, nil, 100, "[]", 0},
	}
	for _, testcase := range testcases {
		actual1, actual2 := GetWeb(testcase.input1, testcase.input2, testcase.input3, testcase.input4)
		if actual1 != testcase.expected1 || actual2 != testcase.expected2 {
			t.Fatalf("utils.GetWeb(\"%v\", %v, %v) result gives\n(\"%v\", %v), "+
				"but not (\"%v\", %v)\n",
				testcase.input1, testcase.input2, testcase.input3,
				actual1, actual2, testcase.expected1, testcase.expected2)
		}
	}
}
