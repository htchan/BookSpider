package model

import (
	"testing"
)

func (w Writer) Equal(compare Writer) bool {
	return w.ID == compare.ID && w.Name == compare.Name
}

func Test_NewWriter(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		writerName string
		expect     Writer
	}{
		{
			name:       "works",
			writerName: "writer",
			expect:     Writer{ID: -1, Name: "writer"},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			result := NewWriter(test.writerName)
			if result != test.expect {
				t.Errorf("got: %v, want: %v", result, test.expect)
			}
		})
	}
}
