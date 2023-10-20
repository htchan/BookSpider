package hjwzw

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVendorService_BookURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		bkID string
		want string
	}{
		{
			name: "int book id",
			bkID: "1234",
			want: "https://tw.hjwzw.com/Book/1234/",
		},
		{
			name: "non int book id",
			bkID: "abcd",
			want: "https://tw.hjwzw.com/Book/abcd/",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			serv := VendorService{}
			got := serv.BookURL(test.bkID)

			assert.Equal(t, test.want, got)
		})
	}
}

func TestVendorService_ChapterListURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		bkID string
		want string
	}{
		{
			name: "int book id",
			bkID: "1234",
			want: "https://tw.hjwzw.com/Book/Chapter/1234/",
		},
		{
			name: "non int book id",
			bkID: "abcd",
			want: "https://tw.hjwzw.com/Book/Chapter/abcd/",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			serv := VendorService{}
			got := serv.ChapterListURL(test.bkID)

			assert.Equal(t, test.want, got)
		})
	}
}

func TestVendorService_ChapterURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		resources []string
		want      string
	}{
		{
			name:      "single full http resource input",
			resources: []string{"http://testing.com"},
			want:      "http://testing.com",
		},
		{
			name:      "single uri resource input with slash",
			resources: []string{"/testing"},
			want:      "https://tw.hjwzw.com/testing",
		},
		{
			name:      "single url resource input without slash",
			resources: []string{"testing"},
			want:      "",
		},
		{
			name:      "zero resources input",
			resources: []string{},
			want:      "",
		},
		{
			name:      "multiple resources input",
			resources: []string{"1234", "abcd"},
			want:      "",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			serv := VendorService{}
			got := serv.ChapterURL(test.resources...)

			assert.Equal(t, test.want, got)
		})
	}
}

func TestVendorService_AvailabilityURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		want string
	}{
		{
			name: "happy flow",
			want: "https://tw.hjwzw.com",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			serv := VendorService{}
			got := serv.AvailabilityURL()

			assert.Equal(t, test.want, got)
		})
	}
}
