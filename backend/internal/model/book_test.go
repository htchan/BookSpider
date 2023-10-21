package model

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_NewBook(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		site   string
		id     int
		expect Book
	}{
		{
			name: "works",
			site: "test",
			id:   100,
			expect: Book{
				Site: "test", ID: 100, HashCode: int(time.Now().Unix()),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			result := NewBook(test.site, test.id)
			assert.Equal(t, result, test.expect)
		})
	}
}

func TestBook_HeaderInfo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		bk     *Book
		expect string
	}{
		{
			name: "happy flow",
			bk: &Book{
				Title:  "title",
				Writer: Writer{Name: "writer"},
			},
			expect: "title\nwriter\n--------------------\n\n",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			result := test.bk.HeaderInfo()

			assert.Equal(t, test.expect, result)
		})
	}
}

func TestBook_MarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		bk        Book
		expect    string
		expectErr bool
	}{
		{
			name: "works",
			bk: Book{
				Site: "test", ID: 1, HashCode: 0,
				Title: "title", Writer: Writer{ID: 1, Name: "writer"}, Type: "type",
				UpdateDate: "date", UpdateChapter: "chapter",
				Status: StatusInProgress, IsDownloaded: true,
				Error: errors.New("error"),
			},
			expect:    `{"site":"test","id":1,"hash_code":"0","title":"title","writer":"writer","type":"type","update_date":"date","update_chapter":"chapter","status":"INPROGRESS","is_downloaded":true,"error":"error"}`,
			expectErr: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			result, err := json.Marshal(test.bk)
			if (err != nil) != test.expectErr {
				t.Errorf("got error: %v, expect error: %v", err, test.expectErr)
			}
			if string(result) != test.expect {
				t.Errorf("got:  %v\nwant: %v", string(result), test.expect)
			}
		})
	}
}

func TestBook_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		bk     *Book
		expect string
	}{
		{
			name: "happy flow without hashcode",
			bk: &Book{
				Site: "test",
				ID:   123,
			},
			expect: "test-123",
		},
		{
			name: "happy flow with hashcode",
			bk: &Book{
				Site:     "test",
				ID:       123,
				HashCode: 999,
			},
			expect: "test-123-999",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			result := test.bk.String()

			assert.Equal(t, test.expect, result)
		})
	}
}
