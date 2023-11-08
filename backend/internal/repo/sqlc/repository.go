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
		Site:           bk.Site,
		ID:             int32(bk.ID),
		Title:          toSqlString(bk.Title),
		WriterID:       toSqlInt(bk.Writer.ID),
		WriterChecksum: toSqlString(bk.Writer.Checksum()),
		Type:           toSqlString(bk.Type),
		UpdateDate:     toSqlString(bk.UpdateDate),
		UpdateChapter:  toSqlString(bk.UpdateChapter),
		Status:         bk.Status.String(),
		IsDownloaded:   bk.IsDownloaded,
		Checksum:       toSqlString(bk.Checksum()),
	})
	if err == nil {
		bk.HashCode = int(result.HashCode)
		return nil
	}

	_, err = r.queries.CreateBookWithHash(r.ctx, sqlc.CreateBookWithHashParams{
		Site:           bk.Site,
		ID:             int32(bk.ID),
		HashCode:       int32(bk.HashCode),
		Title:          toSqlString(bk.Title),
		WriterID:       toSqlInt(bk.Writer.ID),
		WriterChecksum: toSqlString(bk.Writer.Checksum()),
		Type:           toSqlString(bk.Type),
		UpdateDate:     toSqlString(bk.UpdateDate),
		UpdateChapter:  toSqlString(bk.UpdateChapter),
		Status:         bk.Status.String(),
		IsDownloaded:   bk.IsDownloaded,
		Checksum:       toSqlString(bk.Checksum()),
	})
	if err != nil {
		return fmt.Errorf("fail to insert book: %v", err)
	}

	return nil
}

func (r *SqlcRepo) UpdateBook(bk *model.Book) error {
	_, err := r.queries.UpdateBook(r.ctx, sqlc.UpdateBookParams{
		Site:           bk.Site,
		ID:             int32(bk.ID),
		HashCode:       int32(bk.HashCode),
		Title:          toSqlString(bk.Title),
		WriterID:       toSqlInt(bk.Writer.ID),
		WriterChecksum: toSqlString(bk.Writer.Checksum()),
		Type:           toSqlString(bk.Type),
		UpdateDate:     toSqlString(bk.UpdateDate),
		UpdateChapter:  toSqlString(bk.UpdateChapter),
		Status:         bk.Status.String(),
		IsDownloaded:   bk.IsDownloaded,
		Checksum:       toSqlString(bk.Checksum()),
	})
	if err != nil {
		return fmt.Errorf("fail to update book: %w", err)
	}

	return nil
}

func (r *SqlcRepo) FindBookById(id int) (*model.Book, error) {
	result, err := r.queries.GetBookByID(r.ctx, sqlc.GetBookByIDParams{
		Site: r.site,
		ID:   int32(id),
	})
	if err != nil {
		return nil, fmt.Errorf("fail to query book by site id: %w", err)
	}

	var bkErr error
	if result.Data != "" {
		bkErr = fmt.Errorf(result.Data)
	}

	return &model.Book{
		Site:     result.Site,
		ID:       int(result.ID),
		HashCode: int(result.HashCode),
		Title:    result.Title.String,
		Writer: model.Writer{
			ID:   int(result.WriterID.Int32),
			Name: result.Name,
		},
		Type:          result.Type.String,
		UpdateDate:    result.UpdateDate.String,
		UpdateChapter: result.UpdateChapter.String,
		Status:        model.StatusFromString(result.Status),
		IsDownloaded:  result.IsDownloaded,
		Error:         bkErr,
	}, nil
}
func (r *SqlcRepo) FindBookByIdHash(id, hash int) (*model.Book, error) {
	result, err := r.queries.GetBookByIDHash(r.ctx, sqlc.GetBookByIDHashParams{
		Site:     r.site,
		ID:       int32(id),
		HashCode: int32(hash),
	})
	if err != nil {
		return nil, fmt.Errorf("fail to query book by site id: %w", err)
	}

	var bkErr error
	if result.Data != "" {
		bkErr = fmt.Errorf(result.Data)
	}

	return &model.Book{
		Site:     result.Site,
		ID:       int(result.ID),
		HashCode: int(result.HashCode),
		Title:    result.Title.String,
		Writer: model.Writer{
			ID:   int(result.WriterID.Int32),
			Name: result.Name,
		},
		Type:          result.Type.String,
		UpdateDate:    result.UpdateDate.String,
		UpdateChapter: result.UpdateChapter.String,
		Status:        model.StatusFromString(result.Status),
		IsDownloaded:  result.IsDownloaded,
		Error:         bkErr,
	}, nil
}
func (r *SqlcRepo) FindBooksByStatus(Status model.StatusCode) (<-chan model.Book, error) {
	results, err := r.queries.ListBooksByStatus(r.ctx, sqlc.ListBooksByStatusParams{
		Site:   r.site,
		Status: Status.String(),
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
				Site:     results[i].Site,
				ID:       int(results[i].ID),
				HashCode: int(results[i].HashCode),
				Title:    results[i].Title.String,
				Writer: model.Writer{
					ID:   int(results[i].WriterID.Int32),
					Name: results[i].Name,
				},
				Type:          results[i].Type.String,
				UpdateDate:    results[i].UpdateDate.String,
				UpdateChapter: results[i].UpdateChapter.String,
				Status:        model.StatusFromString(results[i].Status),
				IsDownloaded:  results[i].IsDownloaded,
				Error:         bkErr,
			}
		}
		close(bkChan)
	}()

	return bkChan, nil
}
func (r *SqlcRepo) FindAllBooks() (<-chan model.Book, error) {
	results, err := r.queries.ListBooks(r.ctx, r.site)
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
				Site:     results[i].Site,
				ID:       int(results[i].ID),
				HashCode: int(results[i].HashCode),
				Title:    results[i].Title.String,
				Writer: model.Writer{
					ID:   int(results[i].WriterID.Int32),
					Name: results[i].Name,
				},
				Type:          results[i].Type.String,
				UpdateDate:    results[i].UpdateDate.String,
				UpdateChapter: results[i].UpdateChapter.String,
				Status:        model.StatusFromString(results[i].Status),
				IsDownloaded:  results[i].IsDownloaded,
				Error:         bkErr,
			}
		}
		close(bkChan)
	}()

	return bkChan, nil
}
func (r *SqlcRepo) FindBooksForUpdate() (<-chan model.Book, error) {
	results, err := r.queries.ListBooksForUpdate(r.ctx, r.site)
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
				Site:     results[i].Site,
				ID:       int(results[i].ID),
				HashCode: int(results[i].HashCode),
				Title:    results[i].Title.String,
				Writer: model.Writer{
					ID:   int(results[i].WriterID.Int32),
					Name: results[i].Name,
				},
				Type:          results[i].Type.String,
				UpdateDate:    results[i].UpdateDate.String,
				UpdateChapter: results[i].UpdateChapter.String,
				Status:        model.StatusFromString(results[i].Status),
				IsDownloaded:  results[i].IsDownloaded,
				Error:         bkErr,
			}
		}
		close(bkChan)
	}()

	return bkChan, nil
}
func (r *SqlcRepo) FindBooksForDownload() (<-chan model.Book, error) {
	results, err := r.queries.ListBooksForDownload(r.ctx, r.site)
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
				Site:     results[i].Site,
				ID:       int(results[i].ID),
				HashCode: int(results[i].HashCode),
				Title:    results[i].Title.String,
				Writer: model.Writer{
					ID:   int(results[i].WriterID.Int32),
					Name: results[i].Name,
				},
				Type:          results[i].Type.String,
				UpdateDate:    results[i].UpdateDate.String,
				UpdateChapter: results[i].UpdateChapter.String,
				Status:        model.StatusFromString(results[i].Status),
				IsDownloaded:  results[i].IsDownloaded,
				Error:         bkErr,
			}
		}
		close(bkChan)
	}()

	return bkChan, nil
}
func (r *SqlcRepo) FindBooksByTitleWriter(title, writer string, limit, offset int) ([]model.Book, error) {
	results, err := r.queries.ListBooksByTitleWriter(r.ctx, sqlc.ListBooksByTitleWriterParams{
		Site:    r.site,
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
			Site:     results[i].Site,
			ID:       int(results[i].ID),
			HashCode: int(results[i].HashCode),
			Title:    results[i].Title.String,
			Writer: model.Writer{
				ID:   int(results[i].WriterID.Int32),
				Name: results[i].Name,
			},
			Type:          results[i].Type.String,
			UpdateDate:    results[i].UpdateDate.String,
			UpdateChapter: results[i].UpdateChapter.String,
			Status:        model.StatusFromString(results[i].Status),
			IsDownloaded:  results[i].IsDownloaded,
			Error:         bkErr,
		}
	}

	return bks, nil
}
func (r *SqlcRepo) FindBooksByRandom(limit int) ([]model.Book, error) {
	results, err := r.queries.ListRandomBooks(r.ctx, sqlc.ListRandomBooksParams{
		Site:    r.site,
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
			Site:     results[i].Site,
			ID:       int(results[i].ID),
			HashCode: int(results[i].HashCode),
			Title:    results[i].Title.String,
			Writer: model.Writer{
				ID:   int(results[i].WriterID.Int32),
				Name: results[i].Name,
			},
			Type:          results[i].Type.String,
			UpdateDate:    results[i].UpdateDate.String,
			UpdateChapter: results[i].UpdateChapter.String,
			Status:        model.StatusFromString(results[i].Status),
			IsDownloaded:  results[i].IsDownloaded,
			Error:         bkErr,
		}
	}

	return bks, nil
}

func (r *SqlcRepo) FindBookGroupByID(id int) (model.BookGroup, error) {
	results, err := r.queries.GetBookGroupByID(r.ctx, sqlc.GetBookGroupByIDParams{
		Site: r.site,
		ID:   int32(id),
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
			Site:     results[i].Site,
			ID:       int(results[i].ID),
			HashCode: int(results[i].HashCode),
			Title:    results[i].Title.String,
			Writer: model.Writer{
				ID:   int(results[i].WriterID.Int32),
				Name: results[i].Name,
			},
			Type:          results[i].Type.String,
			UpdateDate:    results[i].UpdateDate.String,
			UpdateChapter: results[i].UpdateChapter.String,
			Status:        model.StatusFromString(results[i].Status),
			IsDownloaded:  results[i].IsDownloaded,
			Error:         bkErr,
		}
	}

	return group, nil
}

func (r *SqlcRepo) FindBookGroupByIDHash(id, hashCode int) (model.BookGroup, error) {
	results, err := r.queries.GetBookGroupByIDHash(r.ctx, sqlc.GetBookGroupByIDHashParams{
		Site:     r.site,
		ID:       int32(id),
		HashCode: int32(hashCode),
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
			Site:     results[i].Site,
			ID:       int(results[i].ID),
			HashCode: int(results[i].HashCode),
			Title:    results[i].Title.String,
			Writer: model.Writer{
				ID:   int(results[i].WriterID.Int32),
				Name: results[i].Name,
			},
			Type:          results[i].Type.String,
			UpdateDate:    results[i].UpdateDate.String,
			UpdateChapter: results[i].UpdateChapter.String,
			Status:        model.StatusFromString(results[i].Status),
			IsDownloaded:  results[i].IsDownloaded,
			Error:         bkErr,
		}
	}

	return group, nil
}

func (r *SqlcRepo) UpdateBooksStatus() error {
	return r.queries.UpdateBooksStatus(r.ctx, sqlc.UpdateBooksStatusParams{
		Site:       r.site,
		UpdateDate: toSqlString(strconv.Itoa(time.Now().Year() - 1)),
	})
}

func (r *SqlcRepo) FindAllBookIDs() ([]int, error) {
	result, err := r.queries.FindAllBookIDs(r.ctx, r.site)
	if err != nil {
		return nil, fmt.Errorf("sql failed: %w", err)
	}

	results := make([]int, 0, len(result))
	for _, res := range result {
		results = append(results, int(res))
	}

	return results, nil
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
	bkStat, _ := r.queries.BooksStat(r.ctx, r.site)
	nonErrorBkStat, _ := r.queries.NonErrorBooksStat(r.ctx, r.site)
	errorBkStat, _ := r.queries.ErrorBooksStat(r.ctx, r.site)
	downloadedBkStat, _ := r.queries.DownloadedBooksStat(r.ctx, r.site)
	bkStatusStat, _ := r.queries.BooksStatusStat(r.ctx, r.site)
	writerStat, _ := r.queries.WritersStat(r.ctx, r.site)

	StatusCount := make(map[model.StatusCode]int)
	for i := range bkStatusStat {
		StatusCount[model.StatusFromString(bkStatusStat[i].Status)] = int(bkStatusStat[i].Count)
	}

	return repo.Summary{
		BookCount:       int(bkStat.BookCount),
		UniqueBookCount: int(bkStat.UniqueBookCount),
		MaxBookID:       int(bkStat.MaxBookID.(int64)),
		LatestSuccessID: int(nonErrorBkStat.(int64)),
		ErrorCount:      int(errorBkStat),
		DownloadCount:   int(downloadedBkStat),
		WriterCount:     int(writerStat),
		StatusCount:     StatusCount,
	}
}

func (r *SqlcRepo) Close() error {
	return r.db.Close()
}
