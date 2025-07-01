package service

import (
	"context"
	"testing"

	"github.com/htchan/BookSpider/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestServiceImpl_BookInfo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		bk   *model.Book
		want string
	}{
		{
			name: "happy flow",
			bk:   &model.Book{HashCode: 0},
			want: `{"site":"","id":0,"hash_code":"0","title":"","writer":"","type":"","update_date":"","update_chapter":"","status":"ERROR","is_downloaded":false,"error":""}`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := new(ServiceImpl).BookInfo(context.Background(), test.bk)
			assert.Equal(t, test.want, got)
		})
	}
}
