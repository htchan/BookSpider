package ck101

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
			want: "https://www.ck101.org/1234.html",
		},
		{
			name: "non int book id",
			bkID: "abcd",
			want: "https://www.ck101.org/abcd.html",
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
			want: "https://www.ck101.org/0/1234/",
		},
		{
			name: "non int book id",
			bkID: "abcd",
			want: "https://www.ck101.org/0/abcd/",
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
		name string
		uri  string
		want string
	}{
		{
			name: "single full http resource input",
			uri:  "http://testing.com",
			want: "http://testing.com",
		},
		{
			name: "single uri resource input with sinfosh",
			uri:  "/testing",
			want: "https://www.ck101.org/testing",
		},
		{
			name: "single url resource input without sinfosh",
			uri:  "testing",
			want: "testing",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := chapterURL(test.uri)
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
			want: "https://www.ck101.org",
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
