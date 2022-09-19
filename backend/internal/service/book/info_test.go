package book

import (
	"errors"
	"testing"

	"github.com/htchan/BookSpider/internal/model"
)

func Test_Info(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		bk     model.Book
		expect string
	}{
		{
			name: "works",
			bk: model.Book{
				Site: "site", ID: 1, HashCode: 30,
				Title: "title", Writer: model.Writer{ID: 1, Name: "writer"}, Type: "type",
				UpdateDate: "date", UpdateChapter: "chapter",
				Status: model.End, IsDownloaded: false, Error: errors.New("data"),
			},
			expect: `{"site":"site","id":1,"hash_code":"u","title":"title","writer":"writer","type":"type","update_date":"date","update_chapter":"chapter","status":"END","is_downloaded":false,"error":"data"}`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			result := Info(test.bk)
			if result != test.expect {
				t.Error(result, test.expect)
			}
		})
	}
}
