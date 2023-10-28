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
			expect: StatusError,
		},
		{
			name:   "return error status",
			input:  "Error",
			expect: StatusError,
		},
		{
			name:   "return in progress status",
			input:  "inprogress",
			expect: StatusInProgress,
		},
		{
			name:   "return end status",
			input:  "end",
			expect: StatusEnd,
		},
		{
			name:   "return error status if input unrecognize",
			input:  "unknown",
			expect: StatusError,
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
			status: StatusError,
			expect: "ERROR",
		},
		{
			name:   "in progress status",
			status: StatusInProgress,
			expect: "INPROGRESS",
		},
		{
			name:   "end status",
			status: StatusEnd,
			expect: "END",
		},
		{
			name:   "status code from not defined integer",
			status: StatusCode(10),
			expect: "ERROR",
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
