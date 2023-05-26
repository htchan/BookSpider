package model

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
			assert.Equal(t, test.expect, result)
		})
	}
}

func Test_Writer_Checksum(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		writer Writer
		expect string
	}{
		{
			name: "works",
			writer: Writer{
				Name: "論文",
			},
			expect: "6K665paH",
		},
		{
			name: "works",
			writer: Writer{
				Name: "论 文",
			},
			expect: "6K665paH",
		},
		{
			name: "works",
			writer: Writer{
				Name: strings.Repeat("long string", 20),
			},
			expect: "",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			result := test.writer.Checksum()
			assert.Equal(t, test.expect, result)
		})
	}
}
