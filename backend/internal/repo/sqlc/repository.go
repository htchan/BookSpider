package repo

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/repo"
	"github.com/htchan/BookSpider/internal/sqlc"
	_ "github.com/lib/pq"
)

type SqlcRepo struct {
	site    string
	db      *sql.DB
	ctx     context.Context
	queries *sqlc.Queries
}

var _ repo.Repository = &SqlcRepo{}

func toSqlString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: true}
}

func toSqlInt(i int) sql.NullInt32 {
	return sql.NullInt32{Int32: int32(i), Valid: true}
}

func toSqlBool(b bool) sql.NullBool {
	return sql.NullBool{Bool: b, Valid: true}
}

func NewRepo(site string, db *sql.DB) *SqlcRepo {
	return &SqlcRepo{
		site:    site,
		db:      db,
		ctx:     context.Background(),
		queries: sqlc.New(db),
	}
}

func (r *SqlcRepo) CreateBook(bk *model.Book) error {
	result, err := r.queries.CreateBookWithZeroHash(r.ctx, sqlc.CreateBookWithZeroHashParams{
		Site:           toSqlString(bk.Site),
		ID:             toSqlInt(bk.ID),
		Title:          toSqlString(bk.Title),
		WriterID:       toSqlInt(bk.Writer.ID),
		WriterChecksum: toSqlString(bk.Writer.Checksum()),
		Type:           toSqlString(bk.Type),
		UpdateDate:     toSqlString(bk.UpdateDate),
		UpdateChapter:  toSqlString(bk.UpdateChapter),
		Status:         toSqlString(bk.Status.String()),
		IsDownloaded:   toSqlBool(bk.IsDownloaded),
		Checksum:       toSqlString(bk.Checksum()),
	})
	if err == nil {
		bk.HashCode = int(result.HashCode.Int32)
		return nil
	}

	_, err = r.queries.CreateBookWithHash(r.ctx, sqlc.CreateBookWithHashParams{
		Site:           toSqlString(bk.Site),
		ID:             toSqlInt(bk.ID),
		HashCode:       toSqlInt(bk.HashCode),
		Title:          toSqlString(bk.Title),
		WriterID:       toSqlInt(bk.Writer.ID),
		WriterChecksum: toSqlString(bk.Writer.Checksum()),
		Type:           toSqlString(bk.Type),
		UpdateDate:     toSqlString(bk.UpdateDate),
		UpdateChapter:  toSqlString(bk.UpdateChapter),
		Status:         toSqlString(bk.Status.String()),
		IsDownloaded:   toSqlBool(bk.IsDownloaded),
		Checksum:       toSqlString(bk.Checksum()),
	})
	if err != nil {
		return fmt.Errorf("fail to insert book: %v", err)
	}

	return nil
}

func (r *SqlcRepo) UpdateBook(bk *model.Book) error {
	_, err := r.queries.UpdateBook(r.ctx, sqlc.UpdateBookParams{
		Site:           toSqlString(bk.Site),
		ID:             toSqlInt(bk.ID),
		HashCode:       toSqlInt(bk.HashCode),
		Title:          toSqlString(bk.Title),
		WriterID:       toSqlInt(bk.Writer.ID),
		WriterChecksum: toSqlString(bk.Writer.Checksum()),
		Type:           toSqlString(bk.Type),
		UpdateDate:     toSqlString(bk.UpdateDate),
		UpdateChapter:  toSqlString(bk.UpdateChapter),
		Status:         toSqlString(bk.Status.String()),
		IsDownloaded:   toSqlBool(bk.IsDownloaded),
		Checksum:       toSqlString(bk.Checksum()),
	})
	if err != nil {
		return fmt.Errorf("fail to update book: %w", err)
	}

	return nil
}

func (r *SqlcRepo) FindBookById(id int) (*model.Book, error) {
	result, err := r.queries.GetBookByID(r.ctx, sqlc.GetBookByIDParams{
		Site: toSqlString(r.site),
		ID:   toSqlInt(id),
	})
	if err != nil {
		return nil, fmt.Errorf("fail to query book by site id: %w", err)
	}

	var bkErr error
	if result.Data != "" {
		bkErr = fmt.Errorf(result.Data)
	}

	return &model.Book{
		Site:     result.Site.String,
		ID:       int(result.ID.Int32),
		HashCode: int(result.HashCode.Int32),
		Title:    result.Title.String,
		Writer: model.Writer{
			ID:   int(result.WriterID.Int32),
			Name: result.Name,
		},
		Type:          result.Type.String,
		UpdateDate:    result.UpdateDate.String,
		UpdateChapter: result.UpdateChapter.String,
		Status:        model.StatusFromString(result.Status.String),
		IsDownloaded:  result.IsDownloaded.Bool,
		Error:         bkErr,
	}, nil
}
func (r *SqlcRepo) FindBookByIdHash(id, hash int) (*model.Book, error) {
	result, err := r.queries.GetBookByIDHash(r.ctx, sqlc.GetBookByIDHashParams{
		Site:     toSqlString(r.site),
		ID:       toSqlInt(id),
		HashCode: toSqlInt(hash),
	})
	if err != nil {
		return nil, fmt.Errorf("fail to query book by site id: %w", err)
	}

	var bkErr error
	if result.Data != "" {
		bkErr = fmt.Errorf(result.Data)
	}

	return &model.Book{
		Site:     result.Site.String,
		ID:       int(result.ID.Int32),
		HashCode: int(result.HashCode.Int32),
		Title:    result.Title.String,
		Writer: model.Writer{
			ID:   int(result.WriterID.Int32),
			Name: result.Name,
		},
		Type:          result.Type.String,
		UpdateDate:    result.UpdateDate.String,
		UpdateChapter: result.UpdateChapter.String,
		Status:        model.StatusFromString(result.Status.String),
		IsDownloaded:  result.IsDownloaded.Bool,
		Error:         bkErr,
	}, nil
}
func (r *SqlcRepo) FindBooksByStatus(status model.StatusCode) (<-chan model.Book, error) {
	results, err := r.queries.ListBooksByStatus(r.ctx, sqlc.ListBooksByStatusParams{
		Site:   toSqlString(r.site),
		Status: toSqlString(status.String()),
	})
	if err != nil {
		return nil, fmt.Errorf("fail to query book by site id: %w", err)
	}

	bkChan := make(chan model.Book)

	go func() {
		for i := range results {
			var bkErr error
			if results[i].Data != "" {
				bkErr = fmt.Errorf(results[i].Data)
			}

			bkChan <- model.Book{
				Site:     results[i].Site.String,
				ID:       int(results[i].ID.Int32),
				HashCode: int(results[i].HashCode.Int32),
				Title:    results[i].Title.String,
				Writer: model.Writer{
					ID:   int(results[i].WriterID.Int32),
					Name: results[i].Name,
				},
				Type:          results[i].Type.String,
				UpdateDate:    results[i].UpdateDate.String,
				UpdateChapter: results[i].UpdateChapter.String,
				Status:        model.StatusFromString(results[i].Status.String),
				IsDownloaded:  results[i].IsDownloaded.Bool,
				Error:         bkErr,
			}
		}
		close(bkChan)
	}()

	return bkChan, nil
}
func (r *SqlcRepo) FindAllBooks() (<-chan model.Book, error) {
	results, err := r.queries.ListBooks(r.ctx, toSqlString(r.site))
	if err != nil {
		return nil, fmt.Errorf("fail to query book by site id: %w", err)
	}

	bkChan := make(chan model.Book)

	go func() {
		for i := range results {
			var bkErr error
			if results[i].Data != "" {
				bkErr = fmt.Errorf(results[i].Data)
			}

			bkChan <- model.Book{
				Site:     results[i].Site.String,
				ID:       int(results[i].ID.Int32),
				HashCode: int(results[i].HashCode.Int32),
				Title:    results[i].Title.String,
				Writer: model.Writer{
					ID:   int(results[i].WriterID.Int32),
					Name: results[i].Name,
				},
				Type:          results[i].Type.String,
				UpdateDate:    results[i].UpdateDate.String,
				UpdateChapter: results[i].UpdateChapter.String,
				Status:        model.StatusFromString(results[i].Status.String),
				IsDownloaded:  results[i].IsDownloaded.Bool,
				Error:         bkErr,
			}
		}
		close(bkChan)
	}()

	return bkChan, nil
}
func (r *SqlcRepo) FindBooksForUpdate() (<-chan model.Book, error) {
	results, err := r.queries.ListBooksForUpdate(r.ctx, toSqlString(r.site))
	if err != nil {
		return nil, fmt.Errorf("fail to query book by site id: %w", err)
	}

	bkChan := make(chan model.Book)

	go func() {
		for i := range results {
			var bkErr error
			if results[i].Data != "" {
				bkErr = fmt.Errorf(results[i].Data)
			}

			bkChan <- model.Book{
				Site:     results[i].Site.String,
				ID:       int(results[i].ID.Int32),
				HashCode: int(results[i].HashCode.Int32),
				Title:    results[i].Title.String,
				Writer: model.Writer{
					ID:   int(results[i].WriterID.Int32),
					Name: results[i].Name,
				},
				Type:          results[i].Type.String,
				UpdateDate:    results[i].UpdateDate.String,
				UpdateChapter: results[i].UpdateChapter.String,
				Status:        model.StatusFromString(results[i].Status.String),
				IsDownloaded:  results[i].IsDownloaded.Bool,
				Error:         bkErr,
			}
		}
		close(bkChan)
	}()

	return bkChan, nil
}
func (r *SqlcRepo) FindBooksForDownload() (<-chan model.Book, error) {
	results, err := r.queries.ListBooksForDownload(r.ctx, toSqlString(r.site))
	if err != nil {
		return nil, fmt.Errorf("fail to query book by site id: %w", err)
	}

	bkChan := make(chan model.Book)

	go func() {
		for i := range results {
			var bkErr error
			if results[i].Data != "" {
				bkErr = fmt.Errorf(results[i].Data)
			}

			bkChan <- model.Book{
				Site:     results[i].Site.String,
				ID:       int(results[i].ID.Int32),
				HashCode: int(results[i].HashCode.Int32),
				Title:    results[i].Title.String,
				Writer: model.Writer{
					ID:   int(results[i].WriterID.Int32),
					Name: results[i].Name,
				},
				Type:          results[i].Type.String,
				UpdateDate:    results[i].UpdateDate.String,
				UpdateChapter: results[i].UpdateChapter.String,
				Status:        model.StatusFromString(results[i].Status.String),
				IsDownloaded:  results[i].IsDownloaded.Bool,
				Error:         bkErr,
			}
		}
		close(bkChan)
	}()

	return bkChan, nil
}
func (r *SqlcRepo) FindBooksByTitleWriter(title, writer string, limit, offset int) ([]model.Book, error) {
	results, err := r.queries.ListBooksByTitleWriter(r.ctx, sqlc.ListBooksByTitleWriterParams{
		Site:    toSqlString(r.site),
		Column2: toSqlString(fmt.Sprintf("%%%s%%", title)),
		Column3: toSqlString(fmt.Sprintf("%%%s%%", writer)),
		Limit:   int32(limit), Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("fail to query book by site id: %w", err)
	}

	bks := make([]model.Book, len(results))
	for i := range results {
		var bkErr error
		if results[i].Data != "" {
			bkErr = fmt.Errorf(results[i].Data)
		}

		bks[i] = model.Book{
			Site:     results[i].Site.String,
			ID:       int(results[i].ID.Int32),
			HashCode: int(results[i].HashCode.Int32),
			Title:    results[i].Title.String,
			Writer: model.Writer{
				ID:   int(results[i].WriterID.Int32),
				Name: results[i].Name,
			},
			Type:          results[i].Type.String,
			UpdateDate:    results[i].UpdateDate.String,
			UpdateChapter: results[i].UpdateChapter.String,
			Status:        model.StatusFromString(results[i].Status.String),
			IsDownloaded:  results[i].IsDownloaded.Bool,
			Error:         bkErr,
		}
	}

	return bks, nil
}
func (r *SqlcRepo) FindBooksByRandom(limit int) ([]model.Book, error) {
	results, err := r.queries.ListRandomBooks(r.ctx, sqlc.ListRandomBooksParams{
		Site:    toSqlString(r.site),
		Column2: limit,
	})
	if err != nil {
		return nil, fmt.Errorf("fail to query book by site id: %w", err)
	}

	bks := make([]model.Book, len(results))
	for i := range results {
		var bkErr error
		if results[i].Data != "" {
			bkErr = fmt.Errorf(results[i].Data)
		}

		bks[i] = model.Book{
			Site:     results[i].Site.String,
			ID:       int(results[i].ID.Int32),
			HashCode: int(results[i].HashCode.Int32),
			Title:    results[i].Title.String,
			Writer: model.Writer{
				ID:   int(results[i].WriterID.Int32),
				Name: results[i].Name,
			},
			Type:          results[i].Type.String,
			UpdateDate:    results[i].UpdateDate.String,
			UpdateChapter: results[i].UpdateChapter.String,
			Status:        model.StatusFromString(results[i].Status.String),
			IsDownloaded:  results[i].IsDownloaded.Bool,
			Error:         bkErr,
		}
	}

	return bks, nil
}

func (r *SqlcRepo) FindBookGroupByID(id int) (model.BookGroup, error) {
	results, err := r.queries.GetBookGroupByID(r.ctx, sqlc.GetBookGroupByIDParams{
		Site: toSqlString(r.site),
		ID:   toSqlInt(id),
	})
	if err != nil {
		return nil, fmt.Errorf("fail to get book group by site id: %w", err)
	}

	group := make(model.BookGroup, len(results))
	for i := range results {
		var bkErr error
		if results[i].Data != "" {
			bkErr = fmt.Errorf(results[i].Data)
		}

		group[i] = model.Book{
			Site:     results[i].Site.String,
			ID:       int(results[i].ID.Int32),
			HashCode: int(results[i].HashCode.Int32),
			Title:    results[i].Title.String,
			Writer: model.Writer{
				ID:   int(results[i].WriterID.Int32),
				Name: results[i].Name,
			},
			Type:          results[i].Type.String,
			UpdateDate:    results[i].UpdateDate.String,
			UpdateChapter: results[i].UpdateChapter.String,
			Status:        model.StatusFromString(results[i].Status.String),
			IsDownloaded:  results[i].IsDownloaded.Bool,
			Error:         bkErr,
		}
	}

	return group, nil
}

func (r *SqlcRepo) FindBookGroupByIDHash(id, hashCode int) (model.BookGroup, error) {
	results, err := r.queries.GetBookGroupByIDHash(r.ctx, sqlc.GetBookGroupByIDHashParams{
		Site:     toSqlString(r.site),
		ID:       toSqlInt(id),
		HashCode: toSqlInt(hashCode),
	})
	if err != nil {
		return nil, fmt.Errorf("fail to get book group by site id: %w", err)
	}

	group := make(model.BookGroup, len(results))
	for i := range results {
		var bkErr error
		if results[i].Data != "" {
			bkErr = fmt.Errorf(results[i].Data)
		}

		group[i] = model.Book{
			Site:     results[i].Site.String,
			ID:       int(results[i].ID.Int32),
			HashCode: int(results[i].HashCode.Int32),
			Title:    results[i].Title.String,
			Writer: model.Writer{
				ID:   int(results[i].WriterID.Int32),
				Name: results[i].Name,
			},
			Type:          results[i].Type.String,
			UpdateDate:    results[i].UpdateDate.String,
			UpdateChapter: results[i].UpdateChapter.String,
			Status:        model.StatusFromString(results[i].Status.String),
			IsDownloaded:  results[i].IsDownloaded.Bool,
			Error:         bkErr,
		}
	}

	return group, nil
}

func (r *SqlcRepo) UpdateBooksStatus() error {
	return r.queries.UpdateBooksStatus(r.ctx, sqlc.UpdateBooksStatusParams{
		Site:       toSqlString(r.site),
		UpdateDate: toSqlString(strconv.Itoa(time.Now().Year() - 1)),
	})
}

// writer related
func (r *SqlcRepo) SaveWriter(writer *model.Writer) error {
	result, err := r.queries.CreateWriter(r.ctx, sqlc.CreateWriterParams{
		Name:     toSqlString(writer.Name),
		Checksum: toSqlString(writer.Checksum()),
	})
	if err != nil {
		return fmt.Errorf("fail to save writer: %w", err)
	}

	writer.ID = int(result.ID)

	return nil
}

// error related
func (r *SqlcRepo) SaveError(bk *model.Book, e error) error {
	var err error
	if e == nil {
		_, err = r.queries.DeleteError(r.ctx, sqlc.DeleteErrorParams{
			Site: toSqlString(bk.Site),
			ID:   toSqlInt(bk.ID),
		})
	} else {
		_, err = r.queries.CreateError(r.ctx, sqlc.CreateErrorParams{
			Site: toSqlString(bk.Site),
			ID:   toSqlInt(bk.ID),
			Data: toSqlString(e.Error()),
		})
	}
	if err != nil {
		return fmt.Errorf("fail to save error: %w", err)
	}

	bk.Error = e

	return nil
}

func (r *SqlcRepo) backupBooks(path string) error {
	_, err := r.db.Exec(
		fmt.Sprintf(
			`copy (select * from books where site='%s') to '%s/%s/books_%s.csv' 
			csv header quote as '''' force quote *`,
			r.site, path, r.site, time.Now().Format("2006-01-02"),
		),
	)

	if err != nil {
		return fmt.Errorf("backup books: %w", err)
	}
	return nil
}

func (r *SqlcRepo) backupWriters(path string) error {
	_, err := r.db.Exec(
		fmt.Sprintf(
			`copy (
				select distinct(writers.*) from writers join books on writers.id=books.writer_id 
				where books.site='%s'
			) to '%s/%s/writers_%s.csv' csv header quote as '''' force quote *`,
			r.site, path, r.site, time.Now().Format("2006-01-02"),
		),
	)

	if err != nil {
		return fmt.Errorf("backup writers: %w", err)
	}
	return nil
}

func (r *SqlcRepo) backupErrors(path string) error {
	_, err := r.db.Exec(
		fmt.Sprintf(
			`copy (select * from errors where site='%s') to '%s/%s/errors_%s.csv' 
			csv header quote as '''' force quote *`,
			r.site, path, r.site, time.Now().Format("2006-01-02"),
		),
	)

	if err != nil {
		return fmt.Errorf("backup errors: %w", err)
	}
	return nil
}

func (r *SqlcRepo) Backup(path string) error {
	for _, f := range []func(string) error{r.backupBooks, r.backupWriters, r.backupErrors} {
		err := f(path)
		if err != nil {
			return err
		}
	}
	return nil
}

// database
func (r *SqlcRepo) DBStats() sql.DBStats {
	return r.db.Stats()
}

func (r *SqlcRepo) Stats() repo.Summary {
	bkStat, _ := r.queries.BooksStat(r.ctx, toSqlString(r.site))
	nonErrorBkStat, _ := r.queries.NonErrorBooksStat(r.ctx, toSqlString(r.site))
	errorBkStat, _ := r.queries.ErrorBooksStat(r.ctx, toSqlString(r.site))
	downloadedBkStat, _ := r.queries.DownloadedBooksStat(r.ctx, toSqlString(r.site))
	bkStatusStat, _ := r.queries.BooksStatusStat(r.ctx, toSqlString(r.site))
	writerStat, _ := r.queries.WritersStat(r.ctx, toSqlString(r.site))

	statusCount := make(map[model.StatusCode]int)
	for i := range bkStatusStat {
		statusCount[model.StatusFromString(bkStatusStat[i].Status.String)] = int(bkStatusStat[i].Count)
	}

	return repo.Summary{
		BookCount:       int(bkStat.BookCount),
		UniqueBookCount: int(bkStat.UniqueBookCount),
		MaxBookID:       int(bkStat.MaxBookID.(int64)),
		LatestSuccessID: int(nonErrorBkStat.(int64)),
		ErrorCount:      int(errorBkStat),
		DownloadCount:   int(downloadedBkStat),
		WriterCount:     int(writerStat),
		StatusCount:     statusCount,
	}
}

func (r *SqlcRepo) Close() error {
	return r.db.Close()
}
