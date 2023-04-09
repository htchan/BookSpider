package model

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/siongui/gojianfan"
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

func (bk *Book) HeaderInfo() string {
	return bk.Title + "\n" + bk.Writer.Name + "\n" + CONTENT_SEP + "\n\n"
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

func (bk Book) String() string {
	result := fmt.Sprintf("%v-%v", bk.Site, bk.ID)
	if bk.HashCode > 0 {
		result += fmt.Sprintf("-%v", bk.HashCode)
	}
	return result
}

func simplified(s string) string {
	return gojianfan.T2S(s)
}

func strToShortHex(s string) string {
	b := []byte(s)
	encoded := base64.StdEncoding.EncodeToString(b)
	return encoded
}

func (bk Book) Checksum() string {
	if len(bk.Title) > 100 || len(bk.Writer.Name) > 100 {
		log.Printf("[%v-%v-%v] title: %v; writer: %v is too long", bk.Site, bk.ID, bk.HashCode, bk.Title, bk.Writer.Name)
		return ""
	}

	title := strings.ReplaceAll(bk.Title, " ", "")
	writerName := strings.ReplaceAll(bk.Writer.Name, " ", "")
	return strToShortHex(simplified(fmt.Sprintf("%s-%s", title, writerName)))
}
