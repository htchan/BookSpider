package repo

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/repo"
	_ "github.com/lib/pq"
)

type PsqlRepo struct {
	site string
	db   *sql.DB
}

var _ repo.Repostory = &PsqlRepo{}

func NewRepo(site string, db *sql.DB) *PsqlRepo {
	return &PsqlRepo{site: site, db: db}
}

func (r *PsqlRepo) Migrate() error {
	db := *r.db
	driver, err := postgres.WithInstance(&db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("migrate fail: %w", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file:///migrations",
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("migrate fail: %w", err)
	}
	err = m.Up()
	if err != nil {
		log.Printf("migration: %s", err)
	}
	defer m.Close()
	return nil
}

func (r *PsqlRepo) CreateBook(bk *model.Book) error {
	_, err := r.db.Exec(
		`insert into books 
		(site, id, hash_code, title, writer_id, type, update_date, update_chapter, status, is_downloaded) 
		values 
		($1,$2,0,$3,$4,$5,$6,$7,$8, $9)`,
		bk.Site, bk.ID, bk.Title, bk.Writer.ID, bk.Type,
		bk.UpdateDate, bk.UpdateChapter, bk.Status.String(), bk.IsDownloaded,
	)
	if err == nil {
		bk.HashCode = 0
		return nil
	}
	_, err = r.db.Exec(
		`insert into books 
		(site, id, hash_code, title, writer_id, type, update_date, update_chapter, status, is_downloaded) 
		values 
		($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		bk.Site, bk.ID, bk.HashCode, bk.Title, bk.Writer.ID, bk.Type,
		bk.UpdateDate, bk.UpdateChapter, bk.Status.String(), bk.IsDownloaded,
	)
	if err != nil {
		return fmt.Errorf("fail to insert book: %w", err)
	}

	return nil
}
func (r *PsqlRepo) UpdateBook(bk *model.Book) error {
	_, err := r.db.Exec(
		"update books set title=$4, writer_id=$5, type=$6, update_date=$7, update_chapter=$8, status=$9, is_downloaded=$10 where site=$1 and id=$2 and hash_code=$3",
		bk.Site, bk.ID, bk.HashCode, bk.Title, bk.Writer.ID, bk.Type,
		bk.UpdateDate, bk.UpdateChapter, bk.Status.String(), bk.IsDownloaded,
	)
	if err != nil {
		return fmt.Errorf("fail to update book: %w", err)
	}
	return nil
}

const (
	QueryField = `books.site, books.id, books.hash_code, books.title, 
		books.writer_id, coalesce(writers.name, ''), books.type,
		books.update_date, books.update_chapter, books.status, books.is_downloaded, coalesce(errors.data, '')`
	QueryTable = `books left join writers on books.writer_id=writers.id 
		left join errors on books.site=errors.site and books.id=errors.id`
)

func rowsToBook(rows *sql.Rows) (*model.Book, error) {
	var (
		errStr    string
		statusStr string
	)
	bk := new(model.Book)
	err := rows.Scan(
		&bk.Site, &bk.ID, &bk.HashCode, &bk.Title,
		&bk.Writer.ID, &bk.Writer.Name, &bk.Type,
		&bk.UpdateDate, &bk.UpdateChapter, &statusStr, &bk.IsDownloaded, &errStr,
	)
	if err != nil {
		return nil, fmt.Errorf("fail to query book by site id: %w", err)
	}
	if errStr != "" {
		bk.Error = fmt.Errorf(errStr)
	}
	bk.Status = model.StatusFromString(statusStr)
	return bk, nil
}

func rowsToBookChan(rows *sql.Rows) <-chan model.Book {
	c := make(chan model.Book)
	go func() {
		defer rows.Close()
		for rows.Next() {
			bk, _ := rowsToBook(rows)
			c <- *bk
		}
		close(c)
	}()
	return c
}

func (r *PsqlRepo) FindBookById(id int) (*model.Book, error) {
	rows, err := r.db.Query(
		fmt.Sprintf(
			`select %s from %s where books.site=$1 and books.id=$2 order by hash_code desc`,
			QueryField, QueryTable,
		),
		r.site, id,
	)
	if err != nil {
		return nil, fmt.Errorf("fail to query book by site id: %w", err)
	}
	defer rows.Close()
	if rows.Next() {
		return rowsToBook(rows)
	}
	return nil, fmt.Errorf("fail to query book by site id: %w", repo.BookNotExist)
}
func (r *PsqlRepo) FindBookByIdHash(id, hash int) (*model.Book, error) {
	fmt.Printf(
		`select %s from %s 
		where books.site=$1 and books.id=$2 and books.hash_code=$3 
		`,
		QueryField, QueryTable,
	)
	fmt.Println(r.site, id, hash)
	rows, err := r.db.Query(
		fmt.Sprintf(
			`select %s from %s 
			where books.site=$1 and books.id=$2 and books.hash_code=$3`,
			QueryField, QueryTable,
		),
		r.site, id, hash,
	)
	if err != nil {
		return nil, fmt.Errorf("fail to query book by site id hash: %w", err)
	}
	defer rows.Close()
	if rows.Next() {
		return rowsToBook(rows)
	}
	return nil, fmt.Errorf("fail to query book by site id hash: %w", repo.BookNotExist)
}
func (r *PsqlRepo) FindBooksByStatus(status model.StatusCode) (<-chan model.Book, error) {
	rows, err := r.db.Query(
		fmt.Sprintf(
			`select %s from %s 
			where books.site=$1 and books.status=$2 order by hash_code desc`,
			QueryField, QueryTable,
		),
		r.site, status.String(),
	)
	if err != nil {
		return nil, fmt.Errorf("fail to query book by status: %w", err)
	}
	defer rows.Close()

	return rowsToBookChan(rows), nil
}
func (r *PsqlRepo) FindAllBooks() (<-chan model.Book, error) {
	rows, err := r.db.Query(
		fmt.Sprintf(
			`select %s from %s where books.site=$1 order by books.site, books.id, books.hash_code`,
			QueryField, QueryTable),
		r.site,
	)
	if err != nil {
		return nil, fmt.Errorf("fail to query all books in site: %w", err)
	}

	return rowsToBookChan(rows), nil
}
func (r *PsqlRepo) FindBooksForUpdate() (<-chan model.Book, error) {
	rows, err := r.db.Query(
		fmt.Sprintf(
			`select distinct on (books.site, books.id) %s from %s 
			where books.site=$1 order by books.site, books.id desc, books.hash_code desc`,
			QueryField, QueryTable,
		),
		r.site,
	)
	if err != nil {
		return nil, fmt.Errorf("fail to query books for update: %w", err)
	}

	return rowsToBookChan(rows), nil
}
func (r *PsqlRepo) FindBooksForDownload() (<-chan model.Book, error) {
	rows, err := r.db.Query(
		fmt.Sprintf(
			`select %s from %s
			where books.site=$1 and books.status=$2 and books.is_downloaded=$3
			order by books.update_date desc, books.id desc`,
			QueryField, QueryTable,
		),
		r.site, model.StatusCode(model.End).String(), false,
	)
	if err != nil {
		return nil, fmt.Errorf("fail to query books for update: %w", err)
	}

	return rowsToBookChan(rows), nil
}
func (r *PsqlRepo) FindBooksByTitleWriter(title, writer string, limit, offset int) ([]model.Book, error) {
	if title != "" {
		title = fmt.Sprintf("%%%v%%", title)
	}

	if writer != "" {
		writer = fmt.Sprintf("%%%v%%", writer)
	}

	rows, err := r.db.Query(
		fmt.Sprintf(
			`select %s from %s 
			where books.site=$1 and 
				(books.title like $2 or writers.name like $3) and 
				(books.status != $4)
			order by books.update_date desc, books.id desc limit $5 offset $6`,
			QueryField, QueryTable,
		),
		r.site, title, writer, model.StatusCode(model.Error).String(), limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("fail to query book by title writer: %w", err)
	}

	bks := make([]model.Book, 0)
	for bk := range rowsToBookChan(rows) {
		bks = append(bks, bk)
	}

	return bks, nil
}
func (r *PsqlRepo) FindBooksByRandom(limit int) ([]model.Book, error) {
	rows, err := r.db.Query(
		fmt.Sprintf(
			`select %s from %s
			where books.site=$1 and books.is_downloaded=$2
			order by books.site, books.id desc, books.hash_code desc 
			limit $3 offset RANDOM() * (select count(*) - $3 from books where site=$1 and is_downloaded=$2)`,
			QueryField, QueryTable,
		),
		r.site, true, limit,
	)
	if err != nil {
		return nil, fmt.Errorf("fail to query book by status: %w", err)
	}

	bks := make([]model.Book, 0)
	for bk := range rowsToBookChan(rows) {
		bks = append(bks, bk)
	}

	return bks, nil
}

func generateUpdateStatusCondition(length int) string {
	sqlStmt := "(update_date < '" + strconv.Itoa(time.Now().Year()-1) + "'"

	for i := 0; i < length; i++ {
		sqlStmt += fmt.Sprintf(" or update_chapter like $%d", i+2)
	}

	sqlStmt += fmt.Sprintf(
		") and status = $%d and site=$%d",
		length+2, length+3,
	)

	return sqlStmt
}

func (r *PsqlRepo) UpdateBooksStatus() error {
	matchingKeywords := make([]interface{}, 0, len(repo.ChapterEndKeywords))
	for _, keyword := range repo.ChapterEndKeywords {
		matchingKeywords = append(matchingKeywords, fmt.Sprintf("%%%v%%", keyword))
	}

	conditions := generateUpdateStatusCondition(len(matchingKeywords))

	params := append(
		matchingKeywords, model.StatusCode(model.InProgress).String(), r.site,
	)

	_, err := r.db.Exec(
		fmt.Sprintf("update books set is_downloaded=$1 where %v and is_downloaded=true", conditions),
		append([]interface{}{false}, params...)...,
	)
	if err != nil {
		return fmt.Errorf("update book status failed: %w", err)
	}

	_, err = r.db.Exec(
		fmt.Sprintf("update books set status=$1 where %v", conditions),
		append([]interface{}{model.StatusCode(model.End).String()}, params...)...,
	)
	if err != nil {
		return fmt.Errorf("update book status failed: %w", err)
	}

	return nil
}

// writer related
func (r *PsqlRepo) SaveWriter(writer *model.Writer) error {
	rows, err := r.db.Query(
		`insert into writers (name) values ($1) 
		on conflict (name) do update set name=$1 
		returning id`,
		writer.Name,
	)
	if err != nil {
		return fmt.Errorf("fail to save writer: %w", err)
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&writer.ID)
	} else {
		return fmt.Errorf("fail to query writer: %v", writer.Name)
	}

	return nil
}

// error related
func (r *PsqlRepo) SaveError(bk *model.Book, e error) error {
	var err error
	if e == nil {
		_, err = r.db.Exec(
			"delete from errors where site=$1 and id=$2",
			bk.Site, bk.ID,
		)
	} else {
		_, err = r.db.Exec(
			`insert into errors (site, id, data) values ($1, $2, $3)
			on conflict (site, id)
			do update set data=$3`,
			bk.Site, bk.ID, e.Error(),
		)
	}
	if err != nil {
		return fmt.Errorf("fail to save error: %w", err)
	}
	return nil
}

func (r *PsqlRepo) backupBooks(path string) error {
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

func (r *PsqlRepo) backupWriters(path string) error {
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

func (r *PsqlRepo) backupErrors(path string) error {
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

func (r *PsqlRepo) Backup(path string) error {
	for _, f := range []func(string) error{r.backupBooks, r.backupWriters, r.backupErrors} {
		err := f(path)
		if err != nil {
			return err
		}
	}
	return nil
}

// database
func (r *PsqlRepo) DBStats() sql.DBStats {
	return r.db.Stats()
}

func (r *PsqlRepo) Stats() repo.Summary {
	var summary repo.Summary

	rowsBk, err := r.db.Query("select count(*), count(distinct id), max(id) from books where site=$1", r.site)
	if err == nil {
		defer rowsBk.Close()
		if rowsBk.Next() {
			rowsBk.Scan(&summary.BookCount, &summary.UniqueBookCount, &summary.MaxBookID)
		}
	}

	rowsBkSuccess, err := r.db.Query("select max(id) from books where status<>$1 and site=$2", model.StatusCode(model.Error).String(), r.site)
	if err == nil {
		defer rowsBkSuccess.Close()
		if rowsBkSuccess.Next() {
			rowsBkSuccess.Scan(&summary.LatestSuccessID)
		}
	}

	rowsWt, err := r.db.Query("select count(distinct writers.id) from books join writers on books.writer_id=writers.id where site=$1", r.site)
	if err == nil {
		defer rowsWt.Close()
		if rowsWt.Next() {
			rowsWt.Scan(&summary.WriterCount)
		}
	}

	rowsErr, err := r.db.Query("select count(*) from books where site=$1 and status=$2", r.site, model.StatusCode(model.Error).String())
	if err == nil {
		defer rowsErr.Close()
		if rowsErr.Next() {
			rowsErr.Scan(&summary.ErrorCount)
		}
	}

	rowsDownloaded, err := r.db.Query("select count(*) from books where site=$1 and is_downloaded=$2", r.site, true)
	if err == nil {
		defer rowsDownloaded.Close()
		if rowsDownloaded.Next() {
			rowsDownloaded.Scan(&summary.DownloadCount)
		}
	}

	summary.StatusCount = make(map[model.StatusCode]int)
	rowsStatus, err := r.db.Query("select status, count(*) from books where site=$1 group by status", r.site)
	var (
		statusKey   string
		statusValue int
	)
	if err == nil {
		defer rowsStatus.Close()
		for rowsStatus.Next() {
			rowsStatus.Scan(&statusKey, &statusValue)

			summary.StatusCount[model.StatusFromString(statusKey)] = statusValue
		}
	}

	return summary
}

func (r *PsqlRepo) Close() error {
	return r.db.Close()
}
