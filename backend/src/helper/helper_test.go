package helper

import (
	"testing"
)

func TestSearch(t * testing.T) {
	var testCases = []struct {
		input1, input2, expected string
	} {
		{"abc", "(abc)", "abc"},
		{"abc", "(a.c)", "abc"},
		{"abc", "(.*c)", "abc"},
		{"abc", "(a)", "a"},
		{"abc", "(def)", "error"}}
	success := true
	for _, testCase := range testCases {
		actual := Search(testCase.input1, testCase.input2)
		if actual != testCase.expected {
			t.Errorf("helper.Match(\"%v\", \"%v\") result gives \"%v\", but not \"%v\"\n",
				testCase.input1, testCase.input2, actual, testCase.expected)
			success = false
		}
	}
	if success {
		t.Logf("%v test pass", "helper.Match")
	}
}

func TestSearchAll(t *testing.T) {
	var testCases = []struct {
		input1, input2 string
		expected []string
	} {
		{"abacade", "(a.)", []string{"ab", "ac", "ad"}},}
	success := true
	for _, testCase := range testCases {
		actual := SearchAll(testCase.input1, testCase.input2)
		fail := false
		if len(actual) != len(testCase.expected) {
			fail = true
		}
		for i := range actual {
			if fail || actual[i] != testCase.expected[i] {
				fail = false
				break
			}
		}
		if fail {
			t.Errorf("helper.SearchAll(\"%v\", \"%v\") result gives %v, but not %v\n",
				testCase.input1, testCase.input2, actual, testCase.expected)
			success = false
		}
	}
	if success {
		t.Logf("helper.SearchAll test pass")
	}

}

func TestContains(t *testing.T) {
	var testCases = []struct {
		input1 []int
		input2 int
		expected bool
	} {
		{[]int{1, 2, 3, 4}, 5, false},
		{[]int{}, 1, false},
		{[]int{1}, 1, true},
		{[]int{1, 2, 3, 4}, 3, true}}
	success := true
	for _, testCase := range testCases {
		actual := Contains(testCase.input1, testCase.input2)
		if actual != testCase.expected {
			t.Errorf("helper.Contains(%v, %v) result gives %v, but not %v\n",
				testCase.input1, testCase.input2, actual, testCase.expected)
			success = false
		}
	}
	if success {
		t.Logf("helper.Contains test pass")
	}
}

func TestExist(t *testing.T) {
	var testCases = []struct {
		input1 string
		expected bool
	} {
		{"./helper_test.go", true},
		{"./not_exist.go", false}}
	for _, testCase := range testCases {
		actual := Exists(testCase.input1)
		if actual != testCase.expected {
			t.Errorf("helper.Exist(\"%v\") result gives %v, but not %v\n",
				testCase.input1, actual, testCase.expected)
		}
	}
}