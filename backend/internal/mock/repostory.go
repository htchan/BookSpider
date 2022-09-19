package mock

import (
	"database/sql"

	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/repo"
)

type MockRepostory struct {
	Err error
}

var _ repo.Repostory = &MockRepostory{}

func (m MockRepostory) Migrate() error {
	return m.Err
}

func (m MockRepostory) CreateBook(bk *model.Book) error {
	return m.Err
}

func (m MockRepostory) UpdateBook(bk *model.Book) error {
	return m.Err
}

func (m MockRepostory) FindBookById(id int) (*model.Book, error) {
	return &model.Book{ID: id}, m.Err
}

func (m MockRepostory) FindBookByIdHash(id int, hash int) (*model.Book, error) {
	return &model.Book{ID: id, HashCode: hash}, m.Err
}

func (m MockRepostory) FindBooksByStatus(status model.StatusCode) (<-chan model.Book, error) {
	return nil, m.Err
}

func (m MockRepostory) FindAllBooks() (<-chan model.Book, error) {
	return nil, m.Err
}

func (m MockRepostory) FindBooksForUpdate() (<-chan model.Book, error) {
	return nil, m.Err
}

func (m MockRepostory) FindBooksForDownload() (<-chan model.Book, error) {
	return nil, m.Err
}

func (m MockRepostory) FindBooksByTitleWriter(title, writer string, limit, offset int) ([]model.Book, error) {
	return make([]model.Book, 0), m.Err
}

func (m MockRepostory) FindBooksByRandom(limit, offset int) ([]model.Book, error) {
	return make([]model.Book, 0), m.Err
}

func (m MockRepostory) UpdateBooksStatus() error {
	return m.Err
}
func (m MockRepostory) SaveWriter(*model.Writer) error {
	return m.Err
}

func (m MockRepostory) SaveError(*model.Book, error) error {
	return m.Err
}

func (m MockRepostory) Backup(path string) error {
	return m.Err
}

func (m MockRepostory) DBStats() sql.DBStats {
	return sql.DBStats{}
}

func (m MockRepostory) Stats() repo.Summary {
	return repo.Summary{}
}

func (m MockRepostory) Close() error {
	return m.Err
}
