package router

import (
	"database/sql"

	"github.com/htchan/BookSpider/internal/model"
)

type errResp struct {
	Error string `json:"error"`
}

type booksResp struct {
	Books []bookResp `json:"books"`
}

type bookResp struct {
	Site          string `json:"site"`
	ID            int    `json:"id"`
	HashCode      string `json:"hash_code"`
	Title         string `json:"title"`
	Writer        string `json:"writer"`
	Type          string `json:"type"`
	UpdateDate    string `json:"update_date"`
	UpdateChapter string `json:"update_chapter"`
	Status        string `json:"status"`
	IsDownloaded  bool   `json:"is_downloaded"`
	Error         string `json:"error"`
}

func toBookResp(bk *model.Book) bookResp {
	errString := ""
	if bk.Error != nil {
		errString = bk.Error.Error()
	}

	return bookResp{
		Site: bk.Site, ID: bk.ID, HashCode: bk.FormatHashCode(),
		Title: bk.Title, Writer: bk.Writer.Name, Type: bk.Type,
		UpdateDate: bk.UpdateDate, UpdateChapter: bk.UpdateChapter,
		Status: bk.Status.String(), IsDownloaded: bk.IsDownloaded,
		Error: errString,
	}
}

func toBooksResp(bks []model.Book) booksResp {
	resp := make([]bookResp, 0, len(bks))
	for _, bk := range bks {
		resp = append(resp, toBookResp(&bk))
	}

	return booksResp{Books: resp}
}

type dbStatsResp struct {
	Stats []sql.DBStats `json:"stats"`
}
