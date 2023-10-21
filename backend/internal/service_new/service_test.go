package service

import (
	"context"
	"database/sql"
	"testing"

	circuitbreaker "github.com/htchan/BookSpider/internal/client/v2/circuit_breaker"
	"github.com/htchan/BookSpider/internal/client/v2/retry"
	"github.com/htchan/BookSpider/internal/client/v2/simple"
	config "github.com/htchan/BookSpider/internal/config_new"
	"github.com/htchan/BookSpider/internal/parse/goquery"
	sqlc "github.com/htchan/BookSpider/internal/repo/sqlc"
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
		ctx    context.Context
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
					GoquerySelectorsConfig: config.GoquerySelectorsConfig{
						Title:            config.GoquerySelectorConfig{Selector: "title", Attr: ""},
						Writer:           config.GoquerySelectorConfig{Selector: "writer", Attr: ""},
						BookType:         config.GoquerySelectorConfig{Selector: "type", Attr: ""},
						LastUpdate:       config.GoquerySelectorConfig{Selector: "date", Attr: ""},
						LastChapter:      config.GoquerySelectorConfig{Selector: "chapter", Attr: ""},
						BookChapterURL:   config.GoquerySelectorConfig{Selector: "chapters", Attr: "href"},
						BookChapterTitle: config.GoquerySelectorConfig{Selector: "chapters", Attr: ""},
						ChapterTitle:     config.GoquerySelectorConfig{Selector: "chapter-title", Attr: ""},
						ChapterContent:   config.GoquerySelectorConfig{Selector: "chapter-content", Attr: ""},
					},
				},
				db:     nil,
				weight: nil,
				ctx:    nil,
			},
			expectService: &ServiceImp{
				name: "test service",
				conf: config.SiteConfig{
					GoquerySelectorsConfig: config.GoquerySelectorsConfig{
						Title:            config.GoquerySelectorConfig{Selector: "title", Attr: ""},
						Writer:           config.GoquerySelectorConfig{Selector: "writer", Attr: ""},
						BookType:         config.GoquerySelectorConfig{Selector: "type", Attr: ""},
						LastUpdate:       config.GoquerySelectorConfig{Selector: "date", Attr: ""},
						LastChapter:      config.GoquerySelectorConfig{Selector: "chapter", Attr: ""},
						BookChapterURL:   config.GoquerySelectorConfig{Selector: "chapters", Attr: "href"},
						BookChapterTitle: config.GoquerySelectorConfig{Selector: "chapters", Attr: ""},
						ChapterTitle:     config.GoquerySelectorConfig{Selector: "chapter-title", Attr: ""},
						ChapterContent:   config.GoquerySelectorConfig{Selector: "chapter-content", Attr: ""},
					},
				},
				client: retry.NewClient(
					&retry.RetryClientConfig{},
					circuitbreaker.NewClient(
						&circuitbreaker.CircuitBreakerClientConfig{},
						simple.NewClient(&simple.SimpleClientConfig{}),
					),
				),
				parser: func() *goquery.GoqueryParser {
					parser, _ := goquery.LoadParser(&config.GoquerySelectorsConfig{
						Title:            config.GoquerySelectorConfig{Selector: "title", Attr: ""},
						Writer:           config.GoquerySelectorConfig{Selector: "writer", Attr: ""},
						BookType:         config.GoquerySelectorConfig{Selector: "type", Attr: ""},
						LastUpdate:       config.GoquerySelectorConfig{Selector: "date", Attr: ""},
						LastChapter:      config.GoquerySelectorConfig{Selector: "chapter", Attr: ""},
						BookChapterURL:   config.GoquerySelectorConfig{Selector: "chapters", Attr: "href"},
						BookChapterTitle: config.GoquerySelectorConfig{Selector: "chapters", Attr: ""},
						ChapterTitle:     config.GoquerySelectorConfig{Selector: "chapter-title", Attr: ""},
						ChapterContent:   config.GoquerySelectorConfig{Selector: "chapter-content", Attr: ""},
					})
					return parser
				}(),
				rpo: sqlc.NewRepo("test service", nil),
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
				test.args.db,
				test.args.ctx, test.args.weight,
			)

			assert.Equal(t, test.expectService, serv)
			assert.ErrorIs(t, err, test.expectError)
		})
	}
}
