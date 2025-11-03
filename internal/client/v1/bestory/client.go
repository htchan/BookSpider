package bestory

import (
	"context"

	"github.com/htchan/BookSpider/internal/client/v1"
	"github.com/htchan/BookSpider/internal/config/v3"
	"github.com/htchan/goclient"
)

type bestoryClient struct {
	cli     *goclient.Client
	decoder *client.Decoder
}

var _ client.Client = (*bestoryClient)(nil)

func NewClient(ctx context.Context, conf config.ClientConfig) client.Client {
	// panic("80txt is not available yet")
	return &bestoryClient{
		cli:     newClient(ctx, conf),
		decoder: client.NewDecoder(client.DecodeMethod(conf.DecodeMethod)),
	}
}

func (c *bestoryClient) GetBookInfo(ctx context.Context, bookID string) (*client.BookInfo, error) {
	body, err := c.get(ctx, bookURL(bookID))
	if err != nil {
		return nil, err
	}

	bookInfo, parseErr := parseBook(body)
	if parseErr != nil {
		return nil, parseErr
	}

	return bookInfo, nil
}

func (c *bestoryClient) GetBookChapterList(ctx context.Context, bookID string) (client.ChapterEntryList, error) {
	body, err := c.get(ctx, chapterListURL(bookID))
	if err != nil {
		return nil, err
	}

	chapterList, parseErr := parseChapterList(body)
	if parseErr != nil {
		return nil, parseErr
	}

	return chapterList, nil
}

func (c *bestoryClient) GetChapterContent(ctx context.Context, chapter client.ChapterEntry) (*client.ChapterContent, error) {
	body, err := c.get(ctx, chapter.URL)
	if err != nil {
		return nil, err
	}

	chapterContent, parseErr := parseChapter(body)
	if parseErr != nil {
		return nil, parseErr
	}

	return chapterContent, nil
}

func (c *bestoryClient) Available(ctx context.Context) bool {
	body, err := c.get(ctx, availabilityURL())
	if err != nil {
		return false
	}

	return isAvailable(body)
}
