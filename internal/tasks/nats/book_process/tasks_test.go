package bookprocess

import (
	"context"
	"errors"
	"testing"
	"time"

	mockservice "github.com/htchan/BookSpider/internal/mock/service/v1"
	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/service/v1"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewTaskSet(t *testing.T) {
	tests := []struct {
		name           string
		nc             *nats.Conn
		serv           service.BookService
		availableSites []string
		expect         BookProcessTasks
	}{
		{
			name:           "happy flow/empty services",
			nc:             nil,
			serv:           nil,
			availableSites: []string{},
			expect:         BookProcessTasks{},
		},
		{
			name:           "happy flow/non empty services",
			nc:             nil,
			serv:           nil,
			availableSites: []string{"test"},
			expect: BookProcessTasks{
				&BookProcessTask{
					Site:    "test",
					nc:      nil,
					Service: nil,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewTaskSet(tt.nc, tt.serv, tt.availableSites)
			assert.Equal(t, tt.expect, got)
		})
	}
}

func TestWebsiteUpdateTasks_Publish(t *testing.T) {
	nc, err := nats.Connect(connString)
	assert.NoError(t, err)

	t.Cleanup(func() {
		nc.Close()
	})

	tests := []struct {
		name            string
		sites           []string
		getServ         func(*gomock.Controller) service.BookService
		bk              *model.Book
		expect          []string
		expectErr       error
		expectSubscribe func(t *testing.T, nc *nats.Conn)
	}{
		{
			name:  "happy flow/one supported service found",
			sites: []string{"set_publish_happy_flow_one_supported"},
			getServ: func(ctrl *gomock.Controller) service.BookService {
				serv := mockservice.NewMockBookService(ctrl)
				serv.EXPECT().SupportBook(
					&model.Book{
						Site: "set_publish_happy_flow_one_supported", ID: 1, HashCode: 0,
						Title: "title", Writer: model.Writer{ID: 1, Name: "writer"}, Type: "type",
						UpdateDate: "date", UpdateChapter: "chapter",
						Status: model.StatusInProgress, IsDownloaded: true,
						Error: model.Error{Err: errors.New("error")},
					},
				).Return(true).AnyTimes()

				return serv
			},
			bk: &model.Book{
				Site: "set_publish_happy_flow_one_supported", ID: 1, HashCode: 0,
				Title: "title", Writer: model.Writer{ID: 1, Name: "writer"}, Type: "type",
				UpdateDate: "date", UpdateChapter: "chapter",
				Status: model.StatusInProgress, IsDownloaded: true,
				Error: model.Error{Err: errors.New("error")},
			},
			expect:    []string{"set_publish_happy_flow_one_supported"},
			expectErr: nil,
			expectSubscribe: func(t *testing.T, nc *nats.Conn) {
				var gotMsg *nats.Msg
				sub, err := nc.Subscribe("book_spider.books.process.set_publish_happy_flow_one_supported", func(msg *nats.Msg) {
					gotMsg = msg
					assert.Equal(t, `{"book":{"site":"set_publish_happy_flow_one_supported","id":1,"hash_code":0,"title":"title","type":"type","update_date":"date","update_chapter":"chapter","status":"INPROGRESS","is_downloaded":true,"writer":{"id":1,"name":"writer"},"error":"error"},"trace_id":"00000000000000000000000000000000","span_id":"0000000000000000","trace_flags":0}`, string(msg.Data))
				})
				assert.NoError(t, err)
				time.Sleep(time.Millisecond)
				sub.Unsubscribe()
				assert.NotNil(t, gotMsg, "no message received")
			},
		},
		{
			name:  "error/no supported service",
			sites: []string{"tasks_publish_not_supported_site"},
			getServ: func(c *gomock.Controller) service.BookService {
				serv := mockservice.NewMockBookService(c)
				serv.EXPECT().SupportBook(&model.Book{
					Site: "tasks_publish_not_supported_site", ID: 1, HashCode: 0,
					Title: "title", Writer: model.Writer{ID: 1, Name: "writer"}, Type: "type",
					UpdateDate: "date", UpdateChapter: "chapter",
					Status: model.StatusInProgress, IsDownloaded: true,
					Error: model.Error{Err: errors.New("error")},
				}).Return(false).AnyTimes()

				return serv
			},
			bk: &model.Book{
				Site: "tasks_publish_not_supported_site", ID: 1, HashCode: 0,
				Title: "title", Writer: model.Writer{ID: 1, Name: "writer"}, Type: "type",
				UpdateDate: "date", UpdateChapter: "chapter",
				Status: model.StatusInProgress, IsDownloaded: true,
				Error: model.Error{Err: errors.New("error")},
			},
			expect:          nil,
			expectErr:       ErrNotSupportedBook,
			expectSubscribe: func(t *testing.T, nc *nats.Conn) {},
		},
		{
			name:  "error/mismatch site",
			sites: []string{"test"},
			getServ: func(c *gomock.Controller) service.BookService {
				serv := mockservice.NewMockBookService(c)
				serv.EXPECT().SupportBook(&model.Book{
					Site: "tasks_publish_mismatch_site", ID: 1, HashCode: 0,
					Title: "title", Writer: model.Writer{ID: 1, Name: "writer"}, Type: "type",
					UpdateDate: "date", UpdateChapter: "chapter",
					Status: model.StatusInProgress, IsDownloaded: true,
					Error: model.Error{Err: errors.New("error")},
				}).Return(true).AnyTimes()

				return serv
			},
			bk: &model.Book{
				Site: "tasks_publish_mismatch_site", ID: 1, HashCode: 0,
				Title: "title", Writer: model.Writer{ID: 1, Name: "writer"}, Type: "type",
				UpdateDate: "date", UpdateChapter: "chapter",
				Status: model.StatusInProgress, IsDownloaded: true,
				Error: model.Error{Err: errors.New("error")},
			},
			expect:          nil,
			expectErr:       ErrNotSupportedBook,
			expectSubscribe: func(t *testing.T, nc *nats.Conn) {},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tasks := NewTaskSet(nc, test.getServ(ctrl), test.sites)
			assert.Greater(t, len(tasks), 0)

			go func() {
				result, err := tasks.Publish(context.Background(), test.bk)
				assert.Equal(t, test.expect, result)
				assert.ErrorIs(t, err, test.expectErr)
			}()
			test.expectSubscribe(t, nc)
		})
	}
}
