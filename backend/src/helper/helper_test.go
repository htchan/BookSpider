package helper

import (
	"testing"
	"golang.org/x/text/encoding"
)

func Test_GetWeb(t *testing.T) {
	var testcases = []struct {
		input, expected string
	} {
		{"idk", ""},
		{"http://www.google.com/idk", "404"},
		{"https://api.github.com/repos/htchan/Bookspider/downloads", "[]"},
	}
	for _, testcase := range testcases {
		actual := getWeb(testcase.input)
		if actual != testcase.expected {
			t.Fatalf("helper.getWeb(\"%v\") result gives \"%v\", but not \"%v\"\n",
				testcase.input, actual, testcase.expected)
		}
	}
}
func TestGetWeb(t *testing.T) {
	var testcases = []struct {
		input1 string
		input2 int
		input3 *encoding.Decoder
		expected1 string
		expected2 int
	} {
		{"https://api.github.com/repos/htchan/Bookspider/downloads", 1, nil, "[]", 0},
	}
	for _, testcase := range testcases {
		actual1, actual2 := GetWeb(testcase.input1, testcase.input2, testcase.input3)
		if actual1 != testcase.expected1 || actual2 != testcase.expected2 {
			t.Fatalf("helper.GetWeb(\"%v\", %v, %v) result gives (\"%v\", %v), " +
				"but not (\"%v\", %v)\n",
				testcase.input1, testcase.input2, testcase.input3,
				actual1, actual2, testcase.expected1, testcase.expected2)
		}
	}
}

func TestSearch(t *testing.T) {
	var testcases = []struct {
		input1, input2, expected string
	} {
		{"abc", "(abc)", "abc"},
		{"abc", "(a.c)", "abc"},
		{"abc", "(.*c)", "abc"},
		{"abc", "(a)", "a"},
		{"abc", "(def)", "error"},
	}
	for _, testcase := range testcases {
		actual := Search(testcase.input1, testcase.input2)
		if actual != testcase.expected {
			t.Fatalf("helper.Match(\"%v\", \"%v\") result gives \"%v\", but not \"%v\"\n",
				testcase.input1, testcase.input2, actual, testcase.expected)
		}
	}
}

func TestSearchAll(t *testing.T) {
	var testcases = []struct {
		input1, input2 string
		expected []string
	} {
		{"abacade", "(a.)", []string{"ab", "ac", "ad"}},
		{"abacaade", "(a.)", []string{"ab", "ac", "aa"}},
	}
	for _, testcase := range testcases {
		actual := SearchAll(testcase.input1, testcase.input2)
		if len(actual) != len(testcase.expected) {
			t.Fatalf("helper.SearchAll(\"%v\", \"%v\") result gives %v, but not %v\n",
				testcase.input1, testcase.input2, actual, testcase.expected)
		}
		for i := range actual {
			if actual[i] != testcase.expected[i] {
				t.Fatalf("helper.SearchAll(\"%v\", \"%v\") result gives %v, but not %v\n",
					testcase.input1, testcase.input2, actual, testcase.expected)
			}
		}
	}
}

func TestContains(t *testing.T) {
	var testcases = []struct {
		input1 []int
		input2 int
		expected bool
	} {
		{[]int{1, 2, 3, 4}, 5, false},
		{[]int{}, 1, false},
		{[]int{1}, 1, true},
		{[]int{1, 2, 3, 4}, 3, true},
	}
	for _, testcase := range testcases {
		actual := Contains(testcase.input1, testcase.input2)
		if actual != testcase.expected {
			t.Fatalf("helper.Contains(%v, %v) result gives %v, but not %v\n",
				testcase.input1, testcase.input2, actual, testcase.expected)
		}
	}
}

func TestExist(t *testing.T) {
	var testcases = []struct {
		input1 string
		expected bool
	} {
		{"./helper_test.go", true},
		{"./not_exist.go", false},
	}
	for _, testcase := range testcases {
		actual := Exists(testcase.input1)
		if actual != testcase.expected {
			t.Fatalf("helper.Exist(\"%v\") result gives %v, but not %v\n",
				testcase.input1, actual, testcase.expected)
		}
	}
}