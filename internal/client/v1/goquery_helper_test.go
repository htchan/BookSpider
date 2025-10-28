package client

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
)

func TestParseDoc(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		body   string
		assert func(t *testing.T, got *goquery.Document)
	}{
		{
			name: "happy flow",
			body: "<html><head><title>Test</title></head><body><p>Hello, World!</p></body></html>",
			assert: func(t *testing.T, got *goquery.Document) {
				title := got.Find("title").Text()
				assert.Equal(t, "Test", title)
				paragraph := got.Find("body>p").Text()
				assert.Equal(t, "Hello, World!", paragraph)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := ParseDoc(tt.body)
			tt.assert(t, got)
			assert.NoError(t, err)
		})
	}
}

func TestGetGoqueryContentWithoutChildren(t *testing.T) {
	t.Parallel()

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(
		"<html><head><title>Test</title></head><body>123<span>456<span><p><b>Hello,</b> World!<br/>&nbsp;\u00a0</p></span></span></body></html>",
	))
	assert.NoError(t, err)

	tests := []struct {
		name      string
		selection *goquery.Selection
		want      string
	}{
		{
			name:      "happy flow",
			selection: doc.Find("body"),
			want:      "123",
		},
		{
			name:      "empty selection",
			selection: doc.Find(".nonexistent"),
			want:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := GetGoqueryContentWithoutChildren(tt.selection)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetGoqueryContentWithChildren(t *testing.T) {
	t.Parallel()

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(
		"<html><head><title>Test</title></head><body>123<span>456<span><p><b>Hello,</b> World!<br/>&nbsp;\u00a0</p></span></span></body></html>",
	))
	assert.NoError(t, err)

	tests := []struct {
		name      string
		selection *goquery.Selection
		want      string
	}{
		{
			name:      "happy flow",
			selection: doc.Find("body"),
			want:      "123456Hello, World!",
		},
		{
			name:      "empty selection",
			selection: doc.Find(".nonexistent"),
			want:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := GetGoqueryContentWithChildren(tt.selection)
			assert.Equal(t, tt.want, got)
		})
	}
}
