package model

import "testing"

func Test_StatusFromString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		input  string
		expect StatusCode
	}{
		{
			name:   "return error status case insensitive",
			input:  "ErRoR",
			expect: Error,
		},
		{
			name:   "return error status",
			input:  "Error",
			expect: Error,
		},
		{
			name:   "return in progress status",
			input:  "inprogress",
			expect: InProgress,
		},
		{
			name:   "return end status",
			input:  "end",
			expect: End,
		},
		{
			name:   "return error status if input unrecognize",
			input:  "unknown",
			expect: Error,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			result := StatusFromString(test.input)
			if result != test.expect {
				t.Errorf("got: %v, expect: %v", result, test.expect)
			}
		})
	}
}

func TestStatusCode_String(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		status StatusCode
		expect string
	}{
		{
			name:   "error status",
			status: Error,
			expect: "ERROR",
		},
		{
			name:   "in progress status",
			status: InProgress,
			expect: "INPROGRESS",
		},
		{
			name:   "end status",
			status: End,
			expect: "END",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if test.status.String() != test.expect {
				t.Errorf("got: %v, want: %v", test.status.String(), test.expect)
			}
		})
	}
}
