package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
			assert.Equal(t, test.expect, result)
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
			assert.Equal(t, test.expect, test.status.String())
		})
	}
}

func TestStatusCode_MarshalJSON(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		status      StatusCode
		expect      string
		expectError bool
	}{
		{
			name:   "error status",
			status: StatusError,
			expect: `"ERROR"`,
		},
		{
			name:   "in progress status",
			status: StatusInProgress,
			expect: `"INPROGRESS"`,
		},
		{
			name:   "end status",
			status: StatusEnd,
			expect: `"END"`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			result, err := test.status.MarshalJSON()
			assert.NoError(t, err)
			assert.Equal(t, test.expect, string(result))
		})
	}
}

func TestStatusCode_UnmarshalJSON(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		input       string
		expect      StatusCode
		expectError error
	}{
		{
			name:   "error status",
			input:  `"ERROR"`,
			expect: StatusError,
		},
		{
			name:   "in progress status",
			input:  `"INPROGRESS"`,
			expect: StatusInProgress,
		},
		{
			name:   "end status",
			input:  `"END"`,
			expect: StatusEnd,
		},
		{
			name:        "unmarshal error",
			input:       `"unknown"`,
			expectError: nil,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var result StatusCode
			err := result.UnmarshalJSON([]byte(test.input))
			assert.ErrorIs(t, err, test.expectError)
			assert.Equal(t, test.expect, result)
		})
	}
}
