package bookprocess

import (
	"errors"
	"testing"
	"time"

	mocknats "github.com/htchan/BookSpider/internal/mock/nats"
	mockservice "github.com/htchan/BookSpider/internal/mock/service/v1"
	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/service/v1"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewTask(t *testing.T) {
	tests := []struct {
		name         string
		nc           *nats.Conn
		site         string
		serv         service.BookService
		expectedTask *BookProcessTask
	}{
		{
			name: "assign parameters to correct places",
			nc:   nil,
			site: "test_site",
			serv: nil,
			expectedTask: &BookProcessTask{
				Site:    "test_site",
				nc:      nil,
				Service: nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			task := NewTask(test.site, test.nc, test.serv)

			assert.Equal(t, test.expectedTask, task)
		})
	}
}

func TestWebsiteUpdateTask_subject(t *testing.T) {
	tests := []struct {
		name   string
		site   string
		expect string
	}{
		{
			name:   "subject of site name without .",
			site:   "test_site",
			expect: "book_spider.books.process.test_site",
		},
		{
			name:   "subject of service name with .",
			site:   "example.com",
			expect: "book_spider.books.process.example.com",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			task := NewTask(test.site, nil, nil)

			result := task.subject()
			assert.Equal(t, test.expect, result)
		})
	}
}

func TestWebsiteUpdateTask_Publish(t *testing.T) {
	nc, err := nats.Connect(connString)
	assert.NoError(t, err)
	t.Cleanup(func() {
		nc.Close()
	})

	tests := []struct {
		name            string
		site            string
		bk              *model.Book
		expectSubscribe func(*testing.T, *nats.Conn)
		expectErr       error
	}{
		{
			name: "publish success",
			site: "publish_success",
			bk: &model.Book{
				Site: "test", ID: 1, HashCode: 0,
				Title: "title", Writer: model.Writer{ID: 1, Name: "writer"}, Type: "type",
				UpdateDate: "date", UpdateChapter: "chapter",
				Status: model.StatusInProgress, IsDownloaded: true,
				Error: model.Error{Err: errors.New("error")},
			},
			expectSubscribe: func(t *testing.T, nc *nats.Conn) {
				var gotMsg *nats.Msg
				sub, err := nc.Subscribe("book_spider.books.process.publish_success", func(msg *nats.Msg) {
					gotMsg = msg
					assert.Equal(t, `{"book":{"site":"test","id":1,"hash_code":0,"title":"title","type":"type","update_date":"date","update_chapter":"chapter","status":"INPROGRESS","is_downloaded":true,"writer":{"id":1,"name":"writer"},"error":"error"},"trace_id":"00000000000000000000000000000000","span_id":"0000000000000000","trace_flags":0}`, string(msg.Data))
				})
				assert.NoError(t, err)
				time.Sleep(time.Millisecond)
				sub.Unsubscribe()

				assert.NotNil(t, gotMsg)
			},
			expectErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			task := NewTask(test.site, nc, nil)

			go func() {
				err := task.Publish(t.Context(), test.bk)
				assert.ErrorIs(t, err, test.expectErr)
			}()

			test.expectSubscribe(t, nc)
		})
	}
}

func TestWebsiteUpdateTask_Subscribe(t *testing.T) {
	nc, err := nats.Connect(connString)
	assert.NoError(t, err)
	t.Cleanup(func() {
		nc.Close()
	})

	tests := []struct {
		name      string
		site      string
		getServ   func(*gomock.Controller) service.BookService
		publish   func(*testing.T, *nats.Conn)
		expectErr error
	}{
		{
			name: "happy flow",
			site: "subscribe_happy_flow",
			getServ: func(ctrl *gomock.Controller) service.BookService {
				serv := mockservice.NewMockBookService(ctrl)
				bk := &model.Book{
					Site: "subscribe_happy_flow", ID: 1, HashCode: 0,
					Title: "title", Writer: model.Writer{ID: 1, Name: "writer"}, Type: "type",
					UpdateDate: "date", UpdateChapter: "chapter",
					Status: model.StatusInProgress, IsDownloaded: true,
					Error: model.Error{Err: errors.New("error")},
				}
				serv.EXPECT().SupportBook(bk).Return(true)
				serv.EXPECT().ProcessBook(gomock.Any(), bk).Return(nil)
				return serv
			},
			publish: func(t *testing.T, nc *nats.Conn) {
				err := nc.Publish("book_spider.books.process.subscribe_happy_flow", []byte(`{"book":{"site":"subscribe_happy_flow","id":1,"hash_code":0,"title":"title","type":"type","update_date":"date","update_chapter":"chapter","status":"INPROGRESS","is_downloaded":true,"writer":{"id":1,"name":"writer"},"error":"error"},"trace_id":"01234567890123456789012345678901","span_id":"0123456789012345","trace_flags":1}`))
				assert.NoError(t, err)
			},
			expectErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			task := NewTask(test.site, nc, test.getServ(ctrl))

			ctx, err := task.Subscribe(t.Context())
			assert.ErrorIs(t, err, test.expectErr)

			test.publish(t, nc)
			time.Sleep(100 * time.Millisecond)
			if err == nil {
				defer ctx.Stop()
			}
		})
	}

	t.Run("consume each message once", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		bk := &model.Book{
			Site: "subscribe_each_message_once", ID: 1, HashCode: 0,
			Title: "title", Writer: model.Writer{ID: 1, Name: "writer"}, Type: "type",
			UpdateDate: "date", UpdateChapter: "chapter",
			Status: model.StatusInProgress, IsDownloaded: true,
			Error: model.Error{Err: errors.New("error")},
		}
		bk2 := &model.Book{
			Site: "subscribe_each_message_once", ID: 2, HashCode: 0,
			Title: "title 2", Writer: model.Writer{ID: 2, Name: "writer 2"}, Type: "type 2",
			UpdateDate: "date 2", UpdateChapter: "chapter 2",
			Status: model.StatusInProgress, IsDownloaded: true,
			Error: model.Error{Err: errors.New("error")},
		}
		serv := mockservice.NewMockBookService(ctrl)
		site := "subscribe_each_message_once"
		serv.EXPECT().SupportBook(bk).Return(true).Times(1)
		serv.EXPECT().ProcessBook(gomock.Any(), bk).Return(nil).Times(1)
		serv.EXPECT().SupportBook(bk2).Return(true).Times(1)
		serv.EXPECT().ProcessBook(gomock.Any(), bk2).Return(nil).Times(1)

		task := NewTask(site, nc, serv)

		ctx, err := task.Subscribe(t.Context())
		assert.NoError(t, err)
		err = nc.Publish("book_spider.books.process.subscribe_each_message_once", []byte(`{"book":{"site":"subscribe_each_message_once","id":1,"hash_code":0,"title":"title","type":"type","update_date":"date","update_chapter":"chapter","status":"INPROGRESS","is_downloaded":true,"writer":{"id":1,"name":"writer"},"error":"error"}, "trace_id":"XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX", "span_id":"XXXXXXXXXXXXXXXX", "trace_flags":1}`))
		assert.NoError(t, err)
		time.Sleep(time.Millisecond)
		ctx.Stop()

		ctx2, err := task.Subscribe(t.Context())
		assert.NoError(t, err)
		err = nc.Publish("book_spider.books.process.subscribe_each_message_once", []byte(`{"book":{"site":"subscribe_each_message_once","id":2,"hash_code":0,"title":"title 2","type":"type 2","update_date":"date 2","update_chapter":"chapter 2","status":"INPROGRESS","is_downloaded":true,"writer":{"id":2,"name":"writer 2"},"error":"error"}, "trace_id":"XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX", "span_id":"XXXXXXXXXXXXXXXX", "trace_flags":1}`))
		assert.NoError(t, err)
		time.Sleep(time.Millisecond)

		ctx2.Stop()
	})
}

func TestWebsiteUpdateTask_Validate(t *testing.T) {
	tests := []struct {
		name      string
		site      string
		getServ   func(*gomock.Controller) service.BookService
		params    *BookProcessParams
		expectErr error
	}{
		{
			name: "supported web",
			site: "supported_web",
			getServ: func(c *gomock.Controller) service.BookService {
				serv := mockservice.NewMockBookService(c)
				serv.EXPECT().SupportBook(
					&model.Book{Site: "supported_web"},
				).Return(true)

				return serv
			},
			params: &BookProcessParams{
				Book: model.Book{Site: "supported_web"},
			},
			expectErr: nil,
		},
		{
			name: "unsupported web",
			site: "supported_web",
			getServ: func(c *gomock.Controller) service.BookService {
				serv := mockservice.NewMockBookService(c)
				serv.EXPECT().SupportBook(
					&model.Book{Site: "supported_web"},
				).Return(false)

				return serv
			},
			params: &BookProcessParams{
				Book: model.Book{Site: "supported_web"},
			},
			expectErr: ErrNotSupportedBook,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			task := NewTask(test.site, nil, test.getServ(ctrl))

			err := task.Validate(t.Context(), test.params)

			assert.ErrorIs(t, err, test.expectErr)
		})
	}
}

func TestWebsiteUpdateTask_handler(t *testing.T) {
	tests := []struct {
		name    string
		site    string
		getServ func(*gomock.Controller) service.BookService
		getMsg  func(*gomock.Controller) jetstream.Msg
	}{
		{
			name: "happy flow",
			site: "handle_happy_flow",
			getServ: func(ctrl *gomock.Controller) service.BookService {
				serv := mockservice.NewMockBookService(ctrl)
				bk := &model.Book{
					Site: "handle_happy_flow", ID: 1, HashCode: 0,
					Title: "title", Writer: model.Writer{ID: 1, Name: "writer"}, Type: "type",
					UpdateDate: "date", UpdateChapter: "chapter",
					Status: model.StatusInProgress, IsDownloaded: true,
					Error: model.Error{Err: errors.New("error")},
				}
				serv.EXPECT().SupportBook(bk).Return(true)
				serv.EXPECT().ProcessBook(gomock.Any(), bk).Return(nil)

				return serv
			},
			getMsg: func(ctrl *gomock.Controller) jetstream.Msg {
				msg := mocknats.NewMockNatsMsg(ctrl)
				msg.EXPECT().Data().Return([]byte(`{"book":{"site":"handle_happy_flow","id":1,"hash_code":0,"title":"title","type":"type","update_date":"date","update_chapter":"chapter","status":"INPROGRESS","is_downloaded":true,"writer":{"id":1,"name":"writer"},"error":"error"},"trace_id":"00000000000000000000000000000000","span_id":"0000000000000000","trace_flags":0}`))
				msg.EXPECT().Ack()

				return msg
			},
		},
		{
			name: "error/not supported web",
			site: "handle_not_supported_web",
			getServ: func(ctrl *gomock.Controller) service.BookService {
				serv := mockservice.NewMockBookService(ctrl)
				serv.EXPECT().SupportBook(&model.Book{
					Site: "handle_not_supported_web", ID: 1, HashCode: 0,
					Title: "title", Writer: model.Writer{ID: 1, Name: "writer"}, Type: "type",
					UpdateDate: "date", UpdateChapter: "chapter",
					Status: model.StatusInProgress, IsDownloaded: true,
					Error: model.Error{Err: errors.New("error")},
				}).Return(false)

				return serv
			},
			getMsg: func(ctrl *gomock.Controller) jetstream.Msg {
				msg := mocknats.NewMockNatsMsg(ctrl)
				msg.EXPECT().Data().Return([]byte(`{"book":{"site":"handle_not_supported_web","id":1,"hash_code":0,"title":"title","type":"type","update_date":"date","update_chapter":"chapter","status":"INPROGRESS","is_downloaded":true,"writer":{"id":1,"name":"writer"},"error":"error"},"trace_id":"00000000000000000000000000000000","span_id":"0000000000000000","trace_flags":0}`))
				msg.EXPECT().Ack()

				return msg
			},
		},
		{
			name: "error/non json data",
			site: "handle_non_json_data",
			getServ: func(ctrl *gomock.Controller) service.BookService {
				return mockservice.NewMockBookService(ctrl)
			},
			getMsg: func(ctrl *gomock.Controller) jetstream.Msg {
				msg := mocknats.NewMockNatsMsg(ctrl)
				msg.EXPECT().Data().Return([]byte(`non json data`)).Times(2)
				msg.EXPECT().Ack()

				return msg
			},
		},
		{
			name: "error/fail to ack",
			site: "handle_fail_to_ack",
			getServ: func(ctrl *gomock.Controller) service.BookService {
				serv := mockservice.NewMockBookService(ctrl)
				bk := &model.Book{
					Site: "handle_fail_to_ack", ID: 1, HashCode: 0,
					Title: "title", Writer: model.Writer{ID: 1, Name: "writer"}, Type: "type",
					UpdateDate: "date", UpdateChapter: "chapter",
					Status: model.StatusInProgress, IsDownloaded: true,
					Error: model.Error{Err: errors.New("error")},
				}
				serv.EXPECT().SupportBook(bk).Return(true)
				serv.EXPECT().ProcessBook(gomock.Any(), bk).Return(nil)

				return serv
			},
			getMsg: func(ctrl *gomock.Controller) jetstream.Msg {
				msg := mocknats.NewMockNatsMsg(ctrl)
				msg.EXPECT().Data().Return([]byte(`{"book":{"site":"handle_fail_to_ack","id":1,"hash_code":0,"title":"title","type":"type","update_date":"date","update_chapter":"chapter","status":"INPROGRESS","is_downloaded":true,"writer":{"id":1,"name":"writer"},"error":"error"},"trace_id":"00000000000000000000000000000000","span_id":"0000000000000000","trace_flags":0}`))
				msg.EXPECT().Ack().Return(errors.New("ack error"))

				return msg
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			task := NewTask(test.site, nil, test.getServ(ctrl))

			task.handler(test.getMsg(ctrl))
		})
	}
}
