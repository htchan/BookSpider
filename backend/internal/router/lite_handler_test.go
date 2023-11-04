package router

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	servicemock "github.com/htchan/BookSpider/internal/mock/service/v1"
	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/repo"
	"github.com/htchan/BookSpider/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestGeneralLiteHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		services         map[string]service.Service
		prepareRequest   func(*testing.T) *http.Request
		expectStatusCode int
		expectRes        string
	}{
		{
			name: "happy flow",
			services: map[string]service.Service{
				"test": nil,
			},
			prepareRequest: func(t *testing.T) *http.Request {
				t.Helper()
				req, err := http.NewRequest(http.MethodGet, "/", nil)
				assert.NoError(t, err)
				ctx := context.WithValue(req.Context(), URI_PREFIX_KEY, "/lite/novel")

				return req.WithContext(ctx)
			},
			expectStatusCode: 200,
			expectRes: `<html>
	<head>
		<title>Novel</title>
		<style>
			.site_button {
				display: block;
				text-align: center;
				border-style: solid;
				padding: 1em;
			}
		</style>
	</head>
	<body>
		<h1>Novel</h1>




		
	<div
		class="site_button"
		onclick="location.href='/lite/novel/sites/test'"
	>
		test
	</div>
	<br/>


	</body>
</html>
`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			req := test.prepareRequest(t)
			res := httptest.NewRecorder()

			GeneralLiteHandler(test.services).ServeHTTP(res, req)
			assert.Equal(t, test.expectStatusCode, res.Result().StatusCode)
			assert.Equal(t,
				strings.ReplaceAll(strings.ReplaceAll(test.expectRes, "\t", ""), "  ", ""),
				strings.ReplaceAll(strings.ReplaceAll(res.Body.String(), "\t", ""), "  ", ""),
			)
		})
	}
}

func TestSiteLiteHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		services         map[string]service.Service
		prepareRequest   func(*testing.T, *gomock.Controller) *http.Request
		expectStatusCode int
		expectRes        string
	}{
		{
			name: "happy flow",
			services: map[string]service.Service{
				"test": nil,
			},
			prepareRequest: func(t *testing.T, ctrl *gomock.Controller) *http.Request {
				t.Helper()

				serv := servicemock.NewMockService(ctrl)
				serv.EXPECT().Name().Return("test")
				serv.EXPECT().Stats(gomock.Any()).Return(repo.Summary{})

				req, err := http.NewRequest(http.MethodGet, "/", nil)
				assert.NoError(t, err)
				ctx := context.WithValue(req.Context(), URI_PREFIX_KEY, "/lite/novel")
				ctx = context.WithValue(ctx, SERV_KEY, serv)

				return req.WithContext(ctx)
			},
			expectStatusCode: 200,
			expectRes: `<html>

			<head>
			  <title>Novel - test</title>
			  <style>
			  </style>
			</head>
			
			<body>
			  <h1>test</h1>
			  <div>
			    <p>Book Count: 0</p>
			    <p>Writer Count: 0</p>
			    <p>Unique Book Count: 0</p>
			    <p>Max Book ID: 0</p>
			    <p>Latest Success Book ID: 0</p>
			    
			    <p>DownloadCount: 0</p>
			  </div>
			  <hr/>
			  <h2>Search</h2>
			  <div class="search_panel">
			    <form action="search">
			      <label for="fname">Title:</label><br>
			      <input type="text" id="title" name="title"><br>
			      <label for="lname">Writer:</label><br>
			      <input type="text" id="writer" name="writer"><br>
			      <input type="hidden" id="page" name="page" value="0"><br>
			      <input type="hidden" id="per_page" name="per_page" value="10"><br>
			      <input type="submit" value="Submit">
			    </form>
			    <button onclick="location.href='/lite/novel/random?per_page=10'">Random</button>
			  </div>
			</body>
			
			</html>
`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)

			req := test.prepareRequest(t, ctrl)
			res := httptest.NewRecorder()

			SiteLiteHandlerfunc(res, req)
			assert.Equal(t, test.expectStatusCode, res.Result().StatusCode)
			assert.Equal(t,
				strings.ReplaceAll(strings.ReplaceAll(test.expectRes, "\t", ""), "  ", ""),
				strings.ReplaceAll(strings.ReplaceAll(res.Body.String(), "\t", ""), "  ", ""),
			)
		})
	}
}

func TestSearchLiteHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		services         map[string]service.Service
		prepareRequest   func(*testing.T, *gomock.Controller) *http.Request
		expectStatusCode int
		expectRes        string
	}{
		{
			name: "happy flow",
			services: map[string]service.Service{
				"test": nil,
			},
			prepareRequest: func(t *testing.T, ctrl *gomock.Controller) *http.Request {
				t.Helper()

				serv := servicemock.NewMockService(ctrl)
				serv.EXPECT().Name().Return("test")
				serv.EXPECT().QueryBooks(gomock.Any(), "title", "writer", 10, 0).Return(
					[]model.Book{
						{
							Site: "test", ID: 123, HashCode: 100,
							Title: "title", Writer: model.Writer{Name: "writer"},
							Type: "type", UpdateDate: "date", UpdateChapter: "chapter",
							Status: model.StatusEnd, IsDownloaded: true,
						},
					}, nil,
				)

				req, err := http.NewRequest(http.MethodGet, "/", nil)
				assert.NoError(t, err)
				ctx := context.WithValue(req.Context(), URI_PREFIX_KEY, "/lite/novel")
				ctx = context.WithValue(ctx, SERV_KEY, serv)
				ctx = context.WithValue(ctx, TITLE_KEY, "title")
				ctx = context.WithValue(ctx, WRITER_KEY, "writer")
				ctx = context.WithValue(ctx, LIMIT_KEY, 10)
				ctx = context.WithValue(ctx, OFFSET_KEY, 0)

				return req.WithContext(ctx)
			},
			expectStatusCode: 200,
			expectRes: `<html>

			<head>
			  <title>Novel - test</title>
			  <style>
			    .book-box {
			      border-style: solid;
			      padding-left: 1em;
			      padding-right: 1em;
			    }
			  </style>
			</head>
			
			<body>
			  <h1>test</h1>
			  <div>
			    
			    
			      
			  
			  
			  <div class="book-box" onclick="location.href='/lite/novel/books/123-100/'">
			    <p>title - writer</p>
			    <p>date</p>
			    <p>chapter</p>
			    <p>Downloaded</p>
			  </div>
			
			    
			  </div>
			  <!-- TODO: add next page and last page control -->
			</body>
			
			</html>
`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)

			req := test.prepareRequest(t, ctrl)
			res := httptest.NewRecorder()

			SearchLiteHandler(res, req)
			assert.Equal(t, test.expectStatusCode, res.Result().StatusCode)
			assert.Equal(t,
				strings.ReplaceAll(strings.ReplaceAll(test.expectRes, "\t", ""), "  ", ""),
				strings.ReplaceAll(strings.ReplaceAll(res.Body.String(), "\t", ""), "  ", ""),
			)
		})
	}
}

func TestBookLiteHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		services         map[string]service.Service
		prepareRequest   func(*testing.T, *gomock.Controller) *http.Request
		expectStatusCode int
		expectRes        string
	}{
		{
			name: "happy flow",
			services: map[string]service.Service{
				"test": nil,
			},
			prepareRequest: func(t *testing.T, ctrl *gomock.Controller) *http.Request {
				t.Helper()

				serv := servicemock.NewMockService(ctrl)
				serv.EXPECT().Name().Return("test")

				req, err := http.NewRequest(http.MethodGet, "/", nil)
				assert.NoError(t, err)
				ctx := context.WithValue(req.Context(), URI_PREFIX_KEY, "/lite/novel")
				ctx = context.WithValue(ctx, SERV_KEY, serv)
				ctx = context.WithValue(ctx, BOOK_KEY, &model.Book{
					Site: "test", ID: 123, HashCode: 100,
					Title: "title", Writer: model.Writer{Name: "writer"},
					Type: "type", UpdateDate: "date", UpdateChapter: "chapter",
					Status: model.StatusEnd, IsDownloaded: true,
				})
				ctx = context.WithValue(ctx, BOOK_GROUP_KEY, &model.BookGroup{})

				return req.WithContext(ctx)
			},
			expectStatusCode: 200,
			expectRes: `<html>

			<head>
				<title>Novel - test - title</title>
				<style>
				</style>
			</head>
			
			<body>
				<h1>test</h1>
				<div>
						<p>title - writer</p>
						<p>date</p>
						<p>chapter</p>
						
						<a href="/lite/novel/download">Downloaded</a>
						
				</div>
				<h2>Book Group</h2>
				
				
				</body>
			
			</html>
`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)

			req := test.prepareRequest(t, ctrl)
			res := httptest.NewRecorder()

			BookLiteHandler(res, req)
			assert.Equal(t, test.expectStatusCode, res.Result().StatusCode)
			assert.Equal(t,
				strings.ReplaceAll(strings.ReplaceAll(test.expectRes, "\t", ""), "  ", ""),
				strings.ReplaceAll(strings.ReplaceAll(res.Body.String(), "\t", ""), "  ", ""),
			)
		})
	}
}
