package utils

import (
	"testing"
)

func TestMatch(t *testing.T) {
	var testcases = []struct {
		input1, input2 string
		expected       bool
	}{
		{"abc", "(abc)", false},
	}
	for _, testcase := range testcases {
		actual := Match(testcase.input1, testcase.input2)
		if actual != testcase.expected {
			t.Fatalf("utils.Match(\"%v\", \"%v\") result gives\n\"%v\", but not\n\"%v\"\n",
				testcase.input1, testcase.input2, actual, testcase.expected)
		}
	}

}

func TestSearch(t *testing.T) {
	var testcases = []struct {
		input1, input2, expectString string
		expectErrorExist             bool
	}{
		{"abc", "(abc)", "abc", false},
		{"abc", "(a.c)", "abc", false},
		{"abc", "(.*c)", "abc", false},
		{"abc", "(a)", "a", false},
		{"abc", "(def)", "", true},
	}
	for _, testcase := range testcases {
		actualString, actualError := Search(testcase.input1, testcase.input2)
		if actualString != testcase.expectString || (actualError != nil) != testcase.expectErrorExist {
			t.Fatalf("utils.Search(\"%v\", \"%v\") result gives\n\"%v\", %v, but not\n\"%v\", %v\n",
				testcase.input1, testcase.input2, actualString, actualError,
				testcase.expectString, testcase.expectErrorExist)
		}
	}
}

func TestSearchAll(t *testing.T) {
	var testcases = []struct {
		input1, input2 string
		expected       []string
	}{
		{"abacade", "(a.)", []string{"ab", "ac", "ad"}},
		{"abacaade", "(a.)", []string{"ab", "ac", "aa"}},
	}
	for _, testcase := range testcases {
		actual := SearchAll(testcase.input1, testcase.input2)
		if len(actual) != len(testcase.expected) {
			t.Fatalf("utils.SearchAll(\"%v\", \"%v\") result gives\n%v, but not\n%v\n",
				testcase.input1, testcase.input2, actual, testcase.expected)
		}
		for i := range actual {
			if actual[i] != testcase.expected[i] {
				t.Fatalf("utils.SearchAll(\"%v\", \"%v\") result gives\n%v, but not\n%v\n",
					testcase.input1, testcase.input2, actual, testcase.expected)
			}
		}
	}
}

func TestContains(t *testing.T) {
	var testcases = []struct {
		input1   []int
		input2   int
		expected bool
	}{
		{[]int{1, 2, 3, 4}, 5, false},
		{[]int{}, 1, false},
		{[]int{1}, 1, true},
		{[]int{1, 2, 3, 4}, 3, true},
	}
	for _, testcase := range testcases {
		actual := Contains(testcase.input1, testcase.input2)
		if actual != testcase.expected {
			t.Fatalf("utils.Contains(%v, %v) result gives\n%v, but not\n%v\n",
				testcase.input1, testcase.input2, actual, testcase.expected)
		}
	}
}
