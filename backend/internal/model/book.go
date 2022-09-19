package model

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"
)

type Book struct {
	Site          string
	ID            int
	HashCode      int
	Title         string
	Type          string
	UpdateDate    string
	UpdateChapter string
	Status        StatusCode
	IsDownloaded  bool

	Writer Writer
	Error  error
}

func NewBook(site string, id int) Book {
	return Book{
		Site:     site,
		ID:       id,
		HashCode: GenerateHash(),
	}
}

func GenerateHash() int {
	return int(time.Now().Unix())
}

func (bk Book) MarshalJSON() ([]byte, error) {
	errString := ""
	if bk.Error != nil {
		errString = bk.Error.Error()
	}
	return json.Marshal(&struct {
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
	}{
		Site: bk.Site, ID: bk.ID, HashCode: strconv.FormatInt(int64(bk.HashCode), 36),
		Title: bk.Title, Writer: bk.Writer.Name, Type: bk.Type,
		UpdateDate: bk.UpdateDate, UpdateChapter: bk.UpdateChapter,
		Status: bk.Status.String(), IsDownloaded: bk.IsDownloaded,
		Error: errString,
	})
}

func (bk Book) Equal(compare Book) bool {
	return bk.Site == compare.Site && bk.ID == compare.ID && math.Abs(float64(bk.HashCode-compare.HashCode)) < 1000 &&
		bk.Title == compare.Title && bk.Writer == compare.Writer && bk.Type == compare.Type &&
		bk.UpdateDate == compare.UpdateDate && bk.UpdateChapter == compare.UpdateChapter &&
		bk.Status == compare.Status && bk.IsDownloaded == compare.IsDownloaded &&
		((bk.Error == nil && compare.Error == nil) ||
			((bk.Error != nil && compare.Error != nil) && bk.Error.Error() == compare.Error.Error()))
}

func (bk Book) String() string {
	result := fmt.Sprintf("%v-%v", bk.Site, bk.ID)
	if bk.HashCode > 0 {
		result += fmt.Sprintf("-%v", bk.HashCode)
	}
	return result
}
