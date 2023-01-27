package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/htchan/BookSpider/internal/client"
	config "github.com/htchan/BookSpider/internal/config_new"
	"github.com/htchan/BookSpider/internal/parse/goquery"
	psql "github.com/htchan/BookSpider/internal/repo/psql"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/semaphore"
)

func TestServiceImp_Name(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		serv ServiceImp
		want string
	}{
		{
			name: "happy flow",
			serv: ServiceImp{name: "site name"},
			want: "site name",
		},
		{
			name: "empty site name",
			serv: ServiceImp{name: ""},
			want: "",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := test.serv.Name()

			assert.Equal(t, test.want, got)
		})
	}
}

func Test_LoadService(t *testing.T) {
	t.Parallel()

	type Args struct {
		name   string
		conf   config.SiteConfig
		db     *sql.DB
		weight *semaphore.Weighted
		ctx    *context.Context
	}

	tests := []struct {
		name          string
		args          Args
		expectService Service
		expectError   error
	}{
		{
			name: "happy flow",
			args: Args{
				name: "test service",
				conf: config.SiteConfig{
					GoquerySelectorConfig: config.GoquerySelectorConfig{
						Title:          "title",
						Writer:         "writer",
						BookType:       "type",
						LastUpdate:     "date",
						LastChapter:    "chapter",
						BookChapter:    "chapters",
						ChapterTitle:   "chapter-title",
						ChapterContent: "chapter-content",
					},
				},
				db:     nil,
				weight: nil,
				ctx:    nil,
			},
			expectService: &ServiceImp{
				name: "test service",
				conf: config.SiteConfig{
					GoquerySelectorConfig: config.GoquerySelectorConfig{
						Title:          "title",
						Writer:         "writer",
						BookType:       "type",
						LastUpdate:     "date",
						LastChapter:    "chapter",
						BookChapter:    "chapters",
						ChapterTitle:   "chapter-title",
						ChapterContent: "chapter-content",
					},
				},
				client: client.NewClientV2(&config.SiteConfig{}, nil, nil),
				parser: func() *goquery.GoqueryParser {
					parser, _ := goquery.LoadParser(&config.GoquerySelectorConfig{
						Title:          "title",
						Writer:         "writer",
						BookType:       "type",
						LastUpdate:     "date",
						LastChapter:    "chapter",
						BookChapter:    "chapters",
						ChapterTitle:   "chapter-title",
						ChapterContent: "chapter-content",
					})
					return parser
				}(),
				rpo: psql.NewRepo("test service", nil),
			},
			expectError: nil,
		},
		{
			name: "load parser getting error",
			args: Args{
				name:   "test service",
				conf:   config.SiteConfig{},
				db:     nil,
				weight: nil,
				ctx:    nil,
			},
			expectService: nil,
			expectError:   goquery.ErrBookInfoSelectorEmpty,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			serv, err := LoadService(
				test.args.name, test.args.conf,
				test.args.db, test.args.weight,
				test.args.ctx,
			)

			if !errors.Is(err, test.expectError) {
				t.Errorf("err different:")
				t.Errorf("%v", err)
				t.Errorf("%v", test.expectError)
			}

			assert.Equal(t, test.expectService, serv)

			// if !reflect.DeepEqual(serv, test.expectService) {
			// 	t.Errorf("serv different:")
			// 	t.Errorf("%v", serv)
			// 	t.Errorf("%v", test.expectService)
			// }
		})
	}
}
