package vendor

import (
	"flag"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	leak := flag.Bool("leak", false, "check for memory leaks")
	flag.Parse()

	if *leak {
		goleak.VerifyTestMain(m)
	} else {
		os.Exit(m.Run())
	}
}

func TestCheckChapterEnd(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		chapter string
		want    bool
	}{
		{
			name:    "contains keyword",
			chapter: "it is 全文結",
			want:    true,
		},
		{
			name:    "not contains keyword",
			chapter: "it is not end yet",
			want:    false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := CheckChapterEnd(test.chapter)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestCheckDateEnd(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		date string
		want bool
	}{
		{
			name: "valid format more than 2 yr ago",
			date: "2016-01-02",
			want: true,
		},
		{
			name: "valid format this yr",
			date: strconv.Itoa(time.Now().Year()),
			want: false,
		},
		{
			name: "invalid format",
			date: "invalid",
			want: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := CheckDateEnd(test.date)
			assert.Equal(t, test.want, got)
		})
	}
}
