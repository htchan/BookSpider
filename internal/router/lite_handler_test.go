package router

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	servicemock "github.com/htchan/BookSpider/internal/mock/service/v1"
	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/repo"
	"github.com/htchan/BookSpider/internal/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
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
				ctx := context.WithValue(req.Context(), ContextKeyUriPrefix, "/lite/novel")

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
				ctx := context.WithValue(req.Context(), ContextKeyUriPrefix, "/lite/novel")
				ctx = context.WithValue(ctx, ContextKeyServ, serv)

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
			    <form action="/lite/novel/sites/test/search">
			      <label for="fname">Title:</label><br>
			      <input type="text" id="title" name="title"><br>
			      <label for="lname">Writer:</label><br>
			      <input type="text" id="writer" name="writer"><br>
			      <input type="hidden" id="page" name="page" value="0"><br>
			      <input type="hidden" id="per_page" name="per_page" value="10"><br>
			      <input type="submit" value="Submit">
			    </form>
			    <button onclick="location.href='/lite/novel/sites/test/random?per_page=10'">Random</button>
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
			name: "happy flow without pagination",
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
				ctx := context.WithValue(req.Context(), ContextKeyUriPrefix, "/lite/novel")
				ctx = context.WithValue(ctx, ContextKeySiteName, "test")
				ctx = context.WithValue(ctx, ContextKeyServ, serv)
				ctx = context.WithValue(ctx, ContextKeyTitle, "title")
				ctx = context.WithValue(ctx, ContextKeyWriter, "writer")
				ctx = context.WithValue(ctx, ContextKeyPage, 0)
				ctx = context.WithValue(ctx, ContextKeyPerPage, 10)
				ctx = context.WithValue(ctx, ContextKeyLimit, 10)
				ctx = context.WithValue(ctx, ContextKeyOffset, 0)

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
			      margin: 1em;
				}
				.inline {
				  display: inline-block;
				}
				.tag {
				  display: inline-block;
				  background-color: #f0f0f0;
				  border-radius: 0.5em;
				  padding: 0.2em 0.5em;
				  margin: 0.5em;
				  border: 0.2em solid #000;
				}
			  </style>
			  <style>
				.page-button {
				  display: inline-block;
				  margin: 0em 2%;
				  width: 45%;
				  padding: 1% 0em;
				  text-align: center;
				}
			  </style>
			</head>
			
			<body>
			  <h1>test</h1>
			  <div>
			    
			    
			      
			  
			  
			  <div class="book-box" onclick="location.href='/lite/novel/sites/test/books/123-100/'">
				<p class="inline">title - writer</p>
				<div class="tag">test</div>
				<div class="tag" style="background-color: #00ff00;">Downloaded</div>
			    <p>date</p>
			    <p>chapter</p>
			  </div>
			
			    
			  </div>








			  <div class="pagination">
			  	<div class="page-button"></div>
				<div class="page-button"></div>
			  </div>

			</body>
			
			</html>
`,
		},
		{
			name: "happy flow with pagination",
			services: map[string]service.Service{
				"test": nil,
			},
			prepareRequest: func(t *testing.T, ctrl *gomock.Controller) *http.Request {
				t.Helper()

				serv := servicemock.NewMockService(ctrl)
				serv.EXPECT().Name().Return("test")
				serv.EXPECT().QueryBooks(gomock.Any(), "title", "writer", 1, 5).Return(
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
				ctx := context.WithValue(req.Context(), ContextKeyUriPrefix, "/lite/novel")
				ctx = context.WithValue(ctx, ContextKeySiteName, "test")
				ctx = context.WithValue(ctx, ContextKeyServ, serv)
				ctx = context.WithValue(ctx, ContextKeyTitle, "title")
				ctx = context.WithValue(ctx, ContextKeyWriter, "writer")
				ctx = context.WithValue(ctx, ContextKeyPage, 5)
				ctx = context.WithValue(ctx, ContextKeyPerPage, 1)
				ctx = context.WithValue(ctx, ContextKeyLimit, 1)
				ctx = context.WithValue(ctx, ContextKeyOffset, 5)

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
			      margin: 1em;
				}
				.inline {
				  display: inline-block;
				}
				.tag {
				  display: inline-block;
				  background-color: #f0f0f0;
				  border-radius: 0.5em;
				  padding: 0.2em 0.5em;
				  margin: 0.5em;
				  border: 0.2em solid #000;
				}
			  </style>
			  <style>
				.page-button {
				  display: inline-block;
				  margin: 0em 2%;
				  width: 45%;
				  padding: 1% 0em;
				  text-align: center;
				}
			  </style>
			</head>
			
			<body>
			  <h1>test</h1>
			  <div>
			    
			    
			      
			  
			  
			  <div class="book-box" onclick="location.href='/lite/novel/sites/test/books/123-100/'">
				<p class="inline">title - writer</p>
				<div class="tag">test</div>
				<div class="tag" style="background-color: #00ff00;">Downloaded</div>
			    <p>date</p>
			    <p>chapter</p>
			  </div>
			
			    
			  </div>








			  <div class="pagination">
				<div class="page-button" style="border-style: solid;" onclick="location.href='/lite/novel/sites/test/search?title=title&writer=writer&page=4&per_page=1'">Previous</div>
				<div class="page-button" style="border-style: solid;" onclick="location.href='/lite/novel/sites/test/search?title=title&writer=writer&page=6&per_page=1'">Next</div>
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

			SearchLiteHandler(res, req)
			assert.Equal(t, test.expectStatusCode, res.Result().StatusCode)
			assert.Equal(t,
				strings.ReplaceAll(strings.ReplaceAll(test.expectRes, "\t", ""), "  ", ""),
				strings.ReplaceAll(strings.ReplaceAll(res.Body.String(), "\t", ""), "  ", ""),
			)
		})
	}
}

func TestRandomLiteHandler(t *testing.T) {

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
				serv.EXPECT().RandomBooks(gomock.Any(), 10).Return(
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
				ctx := context.WithValue(req.Context(), ContextKeyUriPrefix, "/lite/novel")
				ctx = context.WithValue(ctx, ContextKeyServ, serv)
				ctx = context.WithValue(ctx, ContextKeyLimit, 10)
				ctx = context.WithValue(ctx, ContextKeyOffset, 0)

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
			      margin: 1em;
				}
				.inline {
				  display: inline-block;
				}
				.tag {
				  display: inline-block;
				  background-color: #f0f0f0;
				  border-radius: 0.5em;
				  padding: 0.2em 0.5em;
				  margin: 0.5em;
				  border: 0.2em solid #000;
				}
			  </style>
			  <style>
				.page-button {
				  display: inline-block;
				  margin: 0em 2%;
				  width: 45%;
				  padding: 1% 0em;
				  text-align: center;
				}
			  </style>
			</head>
			
			<body>
			  <h1>test</h1>
			  <div>
			    
			    
			      
			  
			  
			  <div class="book-box" onclick="location.href='/lite/novel/sites/test/books/123-100/'">
				<p class="inline">title - writer</p>
				<div class="tag">test</div>
				<div class="tag" style="background-color: #00ff00;">Downloaded</div>
			    <p>date</p>
			    <p>chapter</p>
			  </div>
			
			    
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

			RandomLiteHandler(res, req)
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
			name: "happy flow without group",
			services: map[string]service.Service{
				"test": nil,
			},
			prepareRequest: func(t *testing.T, ctrl *gomock.Controller) *http.Request {
				t.Helper()

				req, err := http.NewRequest(http.MethodGet, "/", nil)
				assert.NoError(t, err)
				ctx := context.WithValue(req.Context(), ContextKeyUriPrefix, "/lite/novel")
				ctx = context.WithValue(ctx, ContextKeyBook, &model.Book{
					Site: "test", ID: 123, HashCode: 100,
					Title: "title", Writer: model.Writer{Name: "writer"},
					Type: "type", UpdateDate: "date", UpdateChapter: "chapter",
					Status: model.StatusEnd, IsDownloaded: true,
				})
				ctx = context.WithValue(ctx, ContextKeyBookGroup, &model.BookGroup{})

				return req.WithContext(ctx)
			},
			expectStatusCode: 200,
			expectRes: `<html>

			<head>
				<title>Novel - test - title</title>
				<style>
			      .book-box {
			        border-style: solid;
			        padding-left: 1em;
			        padding-right: 1em;
			        margin: 1em;
				  }
				  .inline {
				    display: inline-block;
				  }
				  .tag {
				    display: inline-block;
				    background-color: #f0f0f0;
				    border-radius: 0.5em;
				    padding: 0.2em 0.5em;
				    margin: 0.5em;
				    border: 0.2em solid #000;
				  }
				</style>
			</head>
			
			<body>
				<h1>test</h1>
				<div>
						<p class="inline">title - writer</p>
						<div class="tag" style="background-color: #00ff00;">Downloaded</div>
						<p>date</p>
						<p>chapter</p>
						
						<a href="/lite/novel/sites/test/books/123-100/download?format=txt">Download TXT</a>
						<a href="/lite/novel/sites/test/books/123-100/download?format=epub">Download EPUB</a>
						
				</div>
				<h2>Book Group</h2>
				
				
				</body>
			
			</html>
`,
		},
		{
			name: "happy flow with group",
			services: map[string]service.Service{
				"test": nil,
			},
			prepareRequest: func(t *testing.T, ctrl *gomock.Controller) *http.Request {
				t.Helper()

				req, err := http.NewRequest(http.MethodGet, "/", nil)
				assert.NoError(t, err)
				ctx := context.WithValue(req.Context(), ContextKeyUriPrefix, "/lite/novel")
				ctx = context.WithValue(ctx, ContextKeyBook, &model.Book{
					Site: "test", ID: 123, HashCode: 100,
					Title: "title", Writer: model.Writer{Name: "writer"},
					Type: "type", UpdateDate: "date", UpdateChapter: "chapter",
					Status: model.StatusEnd, IsDownloaded: true,
				})
				ctx = context.WithValue(ctx, ContextKeyBookGroup, &model.BookGroup{
					{
						Site: "test-2", ID: 123, HashCode: 100,
						Title: "title", Writer: model.Writer{Name: "writer"},
						Type: "type", UpdateDate: "date", UpdateChapter: "chapter",
						Status: model.StatusEnd, IsDownloaded: true,
					},
				})

				return req.WithContext(ctx)
			},
			expectStatusCode: 200,
			expectRes: `<html>

			<head>
				<title>Novel - test - title</title>
				<style>
			      .book-box {
			        border-style: solid;
			        padding-left: 1em;
			        padding-right: 1em;
			        margin: 1em;
				  }
				  .inline {
				    display: inline-block;
				  }
				  .tag {
				    display: inline-block;
				    background-color: #f0f0f0;
				    border-radius: 0.5em;
				    padding: 0.2em 0.5em;
				    margin: 0.5em;
				    border: 0.2em solid #000;
				  }
				</style>
			</head>
			
			<body>
				<h1>test</h1>
				<div>
						<p class="inline">title - writer</p>
						<div class="tag" style="background-color: #00ff00;">Downloaded</div>
						<p>date</p>
						<p>chapter</p>
						
						<a href="/lite/novel/sites/test/books/123-100/download?format=txt">Download TXT</a>
						<a href="/lite/novel/sites/test/books/123-100/download?format=epub">Download EPUB</a>
						
				</div>
				<h2>Book Group</h2>





				<div class="book-box" onclick="location.href='/lite/novel/sites/test-2/books/123-100/'">
				<p class="inline">title - writer</p>
				<div class="tag">test-2</div>
				<div class="tag" style="background-color: #00ff00;">Downloaded</div>
				<p>date</p>
				<p>chapter</p>
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

			BookLiteHandler(res, req)
			assert.Equal(t, test.expectStatusCode, res.Result().StatusCode)
			assert.Equal(t,
				strings.ReplaceAll(strings.ReplaceAll(test.expectRes, "\t", ""), "  ", ""),
				strings.ReplaceAll(strings.ReplaceAll(res.Body.String(), "\t", ""), "  ", ""),
			)
		})
	}
}
