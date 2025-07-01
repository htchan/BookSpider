package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/repo"
	"github.com/htchan/BookSpider/internal/sqlc"
	_ "github.com/lib/pq"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type SqlcRepo struct {
	db      *sql.DB
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

func NewRepo(db *sql.DB) *SqlcRepo {
	return &SqlcRepo{
		db:      db,
		queries: sqlc.New(db),
	}
}

func (r *SqlcRepo) CreateBook(ctx context.Context, bk *model.Book) error {
	createBookCtx, createBookSpan := repo.GetTracer().Start(ctx, "create book")
	defer createBookSpan.End()

	_, createBookWithZeroHashSpan := repo.GetTracer().Start(createBookCtx, "create book with zero hash")
	defer createBookWithZeroHashSpan.End()

	zeroHashParams := sqlc.CreateBookWithZeroHashParams{
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
	}
	zeroHashJsonByte, zeroHashJsonErr := json.Marshal(zeroHashParams)
	if zeroHashJsonErr == nil {
		createBookWithZeroHashSpan.SetAttributes(attribute.String("params", string(zeroHashJsonByte)))
	}
	result, err := r.queries.CreateBookWithZeroHash(ctx, zeroHashParams)
	if err == nil {
		bk.HashCode = int(result.HashCode)
		return nil
	}

	_, createBookWithHashSpan := repo.GetTracer().Start(createBookCtx, "create book with hash")
	defer createBookWithHashSpan.End()

	withHashParams := sqlc.CreateBookWithHashParams{
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
	}
	withHashJsonByte, withHashJsonErr := json.Marshal(withHashParams)
	if withHashJsonErr == nil {
		createBookWithHashSpan.SetAttributes(attribute.String("params", string(withHashJsonByte)))
	}

	_, err = r.queries.CreateBookWithHash(ctx, withHashParams)
	if err != nil {
		return fmt.Errorf("fail to insert book: %v", err)
	}

	return nil
}

func (r *SqlcRepo) UpdateBook(ctx context.Context, bk *model.Book) error {
	_, span := repo.GetTracer().Start(ctx, "update book")
	defer span.End()

	params := sqlc.UpdateBookParams{
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
	}
	jsonByte, jsonErr := json.Marshal(params)
	if jsonErr == nil {
		span.SetAttributes(attribute.String("params", string(jsonByte)))
	}

	_, err := r.queries.UpdateBook(ctx, params)
	if err != nil {
		return fmt.Errorf("fail to update book: %w", err)
	}

	return nil
}

func (r *SqlcRepo) FindBookById(ctx context.Context, site string, id int) (*model.Book, error) {
	_, span := repo.GetTracer().Start(ctx, "find book by id")
	defer span.End()

	span.SetAttributes(attribute.String("site", site), attribute.Int("id", id))

	result, err := r.queries.GetBookByID(ctx, sqlc.GetBookByIDParams{
		Site: site,
		ID:   int32(id),
	})
	if err != nil {
		return nil, fmt.Errorf("fail to query book by site id: %w", err)
	}

	var bkErr error
	if result.Data != "" {
		bkErr = errors.New(result.Data)
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
func (r *SqlcRepo) FindBookByIdHash(ctx context.Context, site string, id, hash int) (*model.Book, error) {
	_, span := repo.GetTracer().Start(ctx, "find book by id hash")
	defer span.End()

	span.SetAttributes(
		attribute.String("site", site),
		attribute.Int("id", id),
		attribute.Int("hash", hash),
	)

	result, err := r.queries.GetBookByIDHash(ctx, sqlc.GetBookByIDHashParams{
		Site:     site,
		ID:       int32(id),
		HashCode: int32(hash),
	})
	if err != nil {
		return nil, fmt.Errorf("fail to query book by site id: %w", err)
	}

	var bkErr error
	if result.Data != "" {
		bkErr = errors.New(result.Data)
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
func (r *SqlcRepo) FindBooksByStatus(ctx context.Context, status model.StatusCode) (<-chan model.Book, error) {
	_, span := repo.GetTracer().Start(ctx, "find books by status")
	defer span.End()

	span.SetAttributes(attribute.String("status", status.String()))

	results, err := r.queries.ListBooksByStatus(ctx, status.String())
	if err != nil {
		return nil, fmt.Errorf("fail to query book by site id: %w", err)
	}

	bkChan := make(chan model.Book)

	go func() {
		for i := range results {
			var bkErr error
			if results[i].Data != "" {
				bkErr = errors.New(results[i].Data)
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
func (r *SqlcRepo) FindAllBooks(ctx context.Context, site string) (<-chan model.Book, error) {
	_, span := repo.GetTracer().Start(ctx, "find all books")
	defer span.End()

	span.SetAttributes(attribute.String("site", site))

	results, err := r.queries.ListBooks(ctx, site)
	if err != nil {
		return nil, fmt.Errorf("fail to query book by site id: %w", err)
	}

	bkChan := make(chan model.Book)

	go func() {
		for i := range results {
			var bkErr error
			if results[i].Data != "" {
				bkErr = errors.New(results[i].Data)
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
func (r *SqlcRepo) FindBooksForUpdate(ctx context.Context, site string) (<-chan model.Book, error) {
	_, span := repo.GetTracer().Start(ctx, "find books for update")
	defer span.End()

	span.SetAttributes(attribute.String("site", site))

	results, err := r.queries.ListBooksForUpdate(ctx, site)
	if err != nil {
		return nil, fmt.Errorf("fail to query book by site id: %w", err)
	}

	bkChan := make(chan model.Book)

	go func() {
		for i := range results {
			var bkErr error
			if results[i].Data != "" {
				bkErr = errors.New(results[i].Data)
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
func (r *SqlcRepo) FindBooksForDownload(ctx context.Context, site string) (<-chan model.Book, error) {
	_, span := repo.GetTracer().Start(ctx, "find books for download")
	defer span.End()

	span.SetAttributes(attribute.String("site", site))

	results, err := r.queries.ListBooksForDownload(ctx, site)
	if err != nil {
		return nil, fmt.Errorf("fail to query book by site id: %w", err)
	}

	bkChan := make(chan model.Book)

	go func() {
		for i := range results {
			var bkErr error
			if results[i].Data != "" {
				bkErr = errors.New(results[i].Data)
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
func (r *SqlcRepo) FindBooksByTitleWriter(ctx context.Context, title, writer string, limit, offset int) ([]model.Book, error) {
	_, span := repo.GetTracer().Start(ctx, "find books by title and writer")
	defer span.End()

	span.SetAttributes(
		attribute.String("title", title),
		attribute.String("writer", writer),
		attribute.Int("limit", limit),
		attribute.Int("offset", offset),
	)

	results, err := r.queries.ListBooksByTitleWriter(ctx, sqlc.ListBooksByTitleWriterParams{
		Column1: toSqlString(fmt.Sprintf("%%%s%%", title)),
		Column2: toSqlString(fmt.Sprintf("%%%s%%", writer)),
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("fail to query book by site id: %w", err)
	}

	bks := make([]model.Book, len(results))
	for i := range results {
		var bkErr error
		if results[i].Data != "" {
			bkErr = errors.New(results[i].Data)
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
func (r *SqlcRepo) FindBooksByRandom(ctx context.Context, limit int) ([]model.Book, error) {
	_, span := repo.GetTracer().Start(ctx, "find books by random")
	defer span.End()

	span.SetAttributes(
		attribute.Int("limit", limit),
	)

	results, err := r.queries.ListRandomBooks(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("fail to query book by site id: %w", err)
	}

	bks := make([]model.Book, len(results))
	for i := range results {
		var bkErr error
		if results[i].Data != "" {
			bkErr = errors.New(results[i].Data)
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

func (r *SqlcRepo) FindBookGroupByID(ctx context.Context, site string, id int) (model.BookGroup, error) {
	_, span := repo.GetTracer().Start(ctx, "find book group by id")
	defer span.End()

	span.SetAttributes(attribute.String("site", site), attribute.Int("id", id))

	results, err := r.queries.GetBookGroupByID(ctx, sqlc.GetBookGroupByIDParams{
		Site: site,
		ID:   int32(id),
	})
	if err != nil {
		return nil, fmt.Errorf("fail to get book group by site id: %w", err)
	}

	group := make(model.BookGroup, len(results))
	for i := range results {
		var bkErr error
		if results[i].Data != "" {
			bkErr = errors.New(results[i].Data)
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

func (r *SqlcRepo) FindBookGroupByIDHash(ctx context.Context, site string, id, hashCode int) (model.BookGroup, error) {
	_, span := repo.GetTracer().Start(ctx, "find book group by id hash")
	defer span.End()

	span.SetAttributes(
		attribute.Int("id", id),
		attribute.Int("hash", hashCode),
	)

	results, err := r.queries.GetBookGroupByIDHash(ctx, sqlc.GetBookGroupByIDHashParams{
		Site:     site,
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
			bkErr = errors.New(results[i].Data)
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

func (r *SqlcRepo) UpdateBooksStatus(ctx context.Context) error {
	_, span := repo.GetTracer().Start(ctx, "update books status")
	defer span.End()

	return r.queries.UpdateBooksStatus(ctx, toSqlString(strconv.Itoa(time.Now().Year()-1)))
}

func (r *SqlcRepo) FindAllBookIDs(ctx context.Context, site string) ([]int, error) {
	_, span := repo.GetTracer().Start(ctx, "find all book ids")
	defer span.End()

	span.SetAttributes(attribute.String("site", site))

	result, err := r.queries.FindAllBookIDs(ctx, site)
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
func (r *SqlcRepo) SaveWriter(ctx context.Context, writer *model.Writer) error {
	_, span := repo.GetTracer().Start(ctx, "save writer")
	defer span.End()

	span.SetAttributes(
		attribute.String("params.writer_name", writer.Name),
		attribute.String("params.writer_checksum", writer.Checksum()),
	)

	result, err := r.queries.CreateWriter(ctx, sqlc.CreateWriterParams{
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
func (r *SqlcRepo) SaveError(ctx context.Context, bk *model.Book, e error) error {
	_, span := repo.GetTracer().Start(ctx, "save error")
	defer span.End()

	span.SetAttributes(
		attribute.String("params.site", bk.Site),
		attribute.Int("params.id", bk.ID),
	)

	var err error
	if e == nil {
		span.SetAttributes(attribute.String("params.error", "nil"))
		_, err = r.queries.DeleteError(ctx, sqlc.DeleteErrorParams{
			Site: toSqlString(bk.Site),
			ID:   toSqlInt(bk.ID),
		})
	} else {
		span.SetAttributes(attribute.String("params.error", e.Error()))
		_, err = r.queries.CreateError(ctx, sqlc.CreateErrorParams{
			Site: toSqlString(bk.Site),
			ID:   toSqlInt(bk.ID),
			Data: toSqlString(e.Error()),
		})
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		return fmt.Errorf("fail to save error: %w", err)
	}

	bk.Error = e

	return nil
}

func (r *SqlcRepo) backupBooks(ctx context.Context, site, path string) error {
	_, span := repo.GetTracer().Start(ctx, "backup books")
	defer span.End()

	span.SetAttributes(attribute.String("path", path))

	_, err := r.db.Exec(
		fmt.Sprintf(
			`copy (select * from books where site='%s') to '%s/%s/books_%s.csv' 
			csv header quote as '''' force quote *`,
			site, path, site, time.Now().Format("2006-01-02"),
		),
	)

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		return fmt.Errorf("backup books: %w", err)
	}
	return nil
}

func (r *SqlcRepo) backupWriters(ctx context.Context, site, path string) error {
	_, span := repo.GetTracer().Start(ctx, "backup writers")
	defer span.End()
	span.SetAttributes(attribute.String("path", path))
	_, err := r.db.Exec(
		fmt.Sprintf(
			`copy (
				select distinct(writers.*) from writers join books on writers.id=books.writer_id 
				where books.site='%s'
			) to '%s/%s/writers_%s.csv' csv header quote as '''' force quote *`,
			site, path, site, time.Now().Format("2006-01-02"),
		),
	)

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		return fmt.Errorf("backup writers: %w", err)
	}
	return nil
}

func (r *SqlcRepo) backupErrors(ctx context.Context, site, path string) error {
	_, span := repo.GetTracer().Start(ctx, "backup errors")
	defer span.End()
	span.SetAttributes(attribute.String("path", path))

	_, err := r.db.Exec(
		fmt.Sprintf(
			`copy (select * from errors where site='%s') to '%s/%s/errors_%s.csv' 
			csv header quote as '''' force quote *`,
			site, path, site, time.Now().Format("2006-01-02"),
		),
	)

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		return fmt.Errorf("backup errors: %w", err)
	}
	return nil
}

func (r *SqlcRepo) Backup(ctx context.Context, site, path string) error {
	_, span := repo.GetTracer().Start(ctx, "backup")
	defer span.End()

	span.SetAttributes(
		attribute.String("path", path),
		attribute.String("site", site),
	)

	for _, f := range []func(context.Context, string, string) error{r.backupBooks, r.backupWriters, r.backupErrors} {
		err := f(ctx, site, path)
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)

			return err
		}
	}
	return nil
}

// database
func (r *SqlcRepo) DBStats(ctx context.Context) sql.DBStats {
	return r.db.Stats()
}

func (r *SqlcRepo) Stats(ctx context.Context, site string) repo.Summary {
	_, bkSpan := repo.GetTracer().Start(ctx, "get books stat")
	bkStat, _ := r.queries.BooksStat(ctx, site)
	bkSpan.End()

	_, nonErrorBkSpan := repo.GetTracer().Start(ctx, "get non error books stat")
	nonErrorBkStat, _ := r.queries.NonErrorBooksStat(ctx, site)
	nonErrorBkSpan.End()

	_, errorBkSpan := repo.GetTracer().Start(ctx, "get error books stat")
	errorBkStat, _ := r.queries.ErrorBooksStat(ctx, site)
	errorBkSpan.End()

	_, downloadedBkSpan := repo.GetTracer().Start(ctx, "get downloaded books stat")
	downloadedBkStat, _ := r.queries.DownloadedBooksStat(ctx, site)
	downloadedBkSpan.End()

	_, writerStatSpan := repo.GetTracer().Start(ctx, "get writers stat")
	bkStatusStat, _ := r.queries.BooksStatusStat(ctx, site)
	writerStatSpan.End()

	_, writerStatSpan = repo.GetTracer().Start(ctx, "get writers stat")
	writerStat, _ := r.queries.WritersStat(ctx, site)
	writerStatSpan.End()

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
