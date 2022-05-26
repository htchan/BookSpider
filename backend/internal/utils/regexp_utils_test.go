package utils

import (
	"testing"
	"fmt"
)

func TestUtils_Regexp(t *testing.T) {
	t.Run("func Match", func(t *testing.T) {
		t.Run("always return false", func(t *testing.T) {
			result := Match("abc", "(abc)")
			if result {
				t.Errorf("Match return true result")
			}
		})
	})

	t.Run("func Search", func(t *testing.T) {
		var testcases = []struct {
			input1, input2, expectString string
			expectErrorExist             bool
		}{
			{"abc", "(abc)", "abc", false},
			{"abc", "(a.c)", "abc", false},
			{"abc", "(.*c)", "abc", false},
			{"abc", "(a)", "a", false},
			{"abc", "(\\w)", "a", false},
			{"abc", "(def)", "", true},
		}
		for i, testcase := range testcases {
			t.Run(fmt.Sprintf("case %v", i), func(t *testing.T) {
				actualString, actualError := Search(testcase.input1, testcase.input2)
				if actualString != testcase.expectString || (actualError != nil) != testcase.expectErrorExist {
					t.Errorf("utils.Search(\"%v\", \"%v\") result gives\n\"%v\", %v, but not\n\"%v\", %v\n",
						testcase.input1, testcase.input2, actualString, actualError,
						testcase.expectString, testcase.expectErrorExist)
				}
			})
		}
	})
	
	t.Run("func SearchAll", func(t *testing.T) {
		var testcases = []struct {
			input1, input2 string
			expected       []string
		}{
			{"abacade", "(a.)", []string{"ab", "ac", "ad"}},
			{"abacaade", "(a.)", []string{"ab", "ac", "aa"}},
		}
		for i, testcase := range testcases {
			t.Run(fmt.Sprintf("case %v", i), func(t *testing.T) {
				t.Parallel()
				actual := SearchAll(testcase.input1, testcase.input2)
				if len(actual) != len(testcase.expected) {
					t.Errorf("utils.SearchAll(\"%v\", \"%v\") result gives\n%v, but not\n%v\n",
						testcase.input1, testcase.input2, actual, testcase.expected)
				}
				for i := range actual {
					if actual[i] != testcase.expected[i] {
						t.Errorf("utils.SearchAll(\"%v\", \"%v\") result gives\n%v, but not\n%v\n",
							testcase.input1, testcase.input2, actual, testcase.expected)
					}
				}
			})
		}
	})

	t.Run("func Contains", func(t *testing.T) {
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
		for i, testcase := range testcases {
			t.Run(fmt.Sprintf("case %v", i), func(t *testing.T) {
				actual := Contains(testcase.input1, testcase.input2)
				if actual != testcase.expected {
					t.Errorf("utils.Contains(%v, %v) result gives\n%v, but not\n%v\n",
						testcase.input1, testcase.input2, actual, testcase.expected)
				}
			})
		}
	})
}
