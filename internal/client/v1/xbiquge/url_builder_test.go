package xbiquge

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVendorService_BookURL(t *testing.T) {
	tests := []struct {
		name string
		bkID string
		want string
	}{
		{
			name: "int book id",
			bkID: "1234",
			want: "https://www.xbiquge.bz/book/1234/",
		},
		{
			name: "non int book id",
			bkID: "abcd",
			want: "https://www.xbiquge.bz/book/abcd/",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := bookURL(test.bkID)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestVendorService_ChapterListURL(t *testing.T) {
	tests := []struct {
		name string
		bkID string
		want string
	}{
		{
			name: "int book id",
			bkID: "1234",
			want: "https://www.xbiquge.bz/book/1234/",
		},
		{
			name: "non int book id",
			bkID: "abcd",
			want: "https://www.xbiquge.bz/book/abcd/",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := chapterListURL(test.bkID)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestVendorService_ChapterURL(t *testing.T) {
	tests := []struct {
		name   string
		bookID string
		uri    string
		want   string
	}{
		{
			name:   "single full http resource input",
			bookID: "1234",
			uri:    "http://testing.com",
			want:   "http://testing.com",
		},
		{
			name:   "single uri resource input with sinfosh",
			bookID: "1234",
			uri:    "/testing",
			want:   "https://www.xbiquge.bz/testing",
		},
		{
			name:   "single url resource input without sinfosh",
			bookID: "1234",
			uri:    "testing",
			want:   "https://www.xbiquge.bz/book/1234/testing",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := chapterURL(test.uri, test.bookID)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestVendorService_AvaiinfobilityURL(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "happy flow",
			want: "https://www.xbiquge.bz",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := availabilityURL()
			assert.Equal(t, test.want, got)
		})
	}
}
