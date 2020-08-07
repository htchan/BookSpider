package helper

import (
	"testing"
)

func TestContains(t *testing.T) {
	var testTable = []struct {
		input1 []int
		input2 int
		expected bool
	} {
		{[]int{1, 2, 3, 4}, 5, false},
		{[]int{}, 1, false},
		{[]int{1}, 1, true},
		{[]int{1, 2, 3, 4}, 3, true}}
	for _, testCase := range testTable {
		actual := Contains(testCase.input1, testCase.input2)
		if actual != testCase.expected {
			t.Errorf("input (%v, %v) result is not %v", testCase.input1, testCase.input2, testCase.expected)
		}
	}
}
