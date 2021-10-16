package books

import (
	"encoding/json"
	"github.com/htchan/BookSpider/internal/utils"
	"golang.org/x/text/encoding"
	"io/ioutil"
	"log"
	"strconv"

	"errors"
	"fmt"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

const BOOK_MAX_THREAD = 1000

type MetaInfo struct {
	baseUrl, downloadUrl, chapterUrl, chapterUrlPattern                   string
	titleRegex, writerRegex, typeRegex, lastUpdateRegex, lastChapterRegex string
	chapterUrlRegex, chapterTitleRegex, chapterContentRegex               string
}

func NewMetaInfo(info map[string]string) (*MetaInfo, error) {
	metaInfo := new(MetaInfo)
	expectKey := [12]string{
		"baseUrl", "downloadUrl", "chapterUrl", "chapterUrlPattern",
		"titleRegex", "writerRegex", "typeRegex", "lastUpdateRegex",
		"lastChapterRegex", "chapterUrlRegex", "chapterTitleRegex", "chapterContentRegex"}
	for _, key := range expectKey {
		_, ok := info[key]
		if !ok {
			return nil, errors.New("missing key " + key)
		}
	}
	metaInfo.baseUrl = info["baseUrl"]
	metaInfo.downloadUrl = info["downloadUrl"]
	metaInfo.chapterUrl = info["chapterUrl"]
	metaInfo.chapterUrlPattern = info["chapterUrlPattern"]
	metaInfo.titleRegex = info["titleRegex"]
	metaInfo.writerRegex = info["writerRegex"]
	metaInfo.typeRegex = info["typeRegex"]
	metaInfo.lastUpdateRegex = info["lastUpdateRegex"]
	metaInfo.lastChapterRegex = info["lastChapterRegex"]
	metaInfo.chapterUrlRegex = info["chapterUrlRegex"]
	metaInfo.chapterTitleRegex = info["chapterTitleRegex"]
	metaInfo.chapterContentRegex = info["chapterContentRegex"]
	return metaInfo, nil
}

type Book struct {
	SiteName                                     string
	Id, Version                                  int
	Title, Writer, Type, LastUpdate, LastChapter string
	EndFlag, DownloadFlag, ReadFlag              bool
	decoder                                      *encoding.Decoder
	metaInfo                                     MetaInfo
	CONST_SLEEP                                  int
}

func NewBook(siteName string, id int, metaInfo MetaInfo,
	decoder *encoding.Decoder, tx *sql.Tx) *Book {
	book := new(Book)
	book.SiteName = siteName
	book.Id = id
	book.Version = -1
	book.decoder = decoder
	book.metaInfo = metaInfo
	book.metaInfo.baseUrl = fmt.Sprintf(metaInfo.baseUrl, id)
	book.metaInfo.downloadUrl = fmt.Sprintf(metaInfo.downloadUrl, id)
	var rows *sql.Rows
	var err error
	for i := 0; i < 2; i++ {
		rows, err = tx.Query("select max(version) from books where "+
			"site=? and num=?", siteName, id)
		if err != nil {
			if i == 1 {
				panic(err)
			}
			continue
		} else if rows.Next() {
			rows.Scan(&book.Version)
		}
	}
	rows.Close()
	return book
}

func LoadBook(rows *sql.Rows, metaInfo MetaInfo, decoder *encoding.Decoder, constSleep int) (*Book, error) {
	book := new(Book)
	book.decoder = decoder
	book.metaInfo = metaInfo
	book.CONST_SLEEP = constSleep
	err := rows.Scan(&book.SiteName, &book.Id, &book.Version, &book.Title, &book.Writer, &book.Type,
		&book.LastUpdate, &book.LastChapter, &book.EndFlag, &book.DownloadFlag, &book.ReadFlag)
	book.metaInfo.baseUrl = fmt.Sprintf(metaInfo.baseUrl, book.Id)
	book.metaInfo.downloadUrl = fmt.Sprintf(metaInfo.downloadUrl, book.Id)
	return book, err
}

func (book Book) Log(info map[string]interface{}) {
	info["site"], info["id"], info["version"] = book.SiteName, book.Id, book.Version
	outputByte, err := json.Marshal(info)
	utils.CheckError(err)
	log.Println(string(outputByte))
}

func (book *Book) validHTML(html string, url string, trial int) bool {
	if len(html) == 0 {
		book.Log(map[string]interface{}{
			"retry": trial, "url": url, "message": "load html fail - zero length",
		})
		return false
	} else if _, err := strconv.Atoi(html); err == nil {
		book.Log(map[string]interface{}{
			"retry": trial, "url": url, "message": "load html fail - code " + html,
		})
		return false
	} else {
		book.Log(map[string]interface{}{
			"retry": trial, "url": url, "message": "load html success",
		})
	}
	return true
}

func (book Book) Content(bookStoragePath string) string {
	if book.Title == "" || !book.DownloadFlag {
		return ""
	}
	bookLocation := book.StorageLocation(bookStoragePath)
	content, err := ioutil.ReadFile(bookLocation)
	utils.CheckError(err)
	return string(content)
}

func (book Book) StorageLocation(storagePath string) (bookLocation string) {
	bookLocation = storagePath + "/" + strconv.Itoa(book.Id)
	if book.Version > 0 {
		bookLocation += "-v" + strconv.Itoa(book.Version)
	}
	bookLocation += ".txt"
	return
}

// to string function
func (book Book) String() string {
	return book.SiteName + "\t" + strconv.Itoa(book.Id) + "\t" + strconv.Itoa(book.Version) + "\n" +
		book.Title + "\t" + book.Writer + "\n" + book.LastUpdate + "\t" + book.LastChapter
}

func (book Book) Map() map[string]interface{} {
	return map[string]interface{}{
		"site":     book.SiteName,
		"id":       book.Id,
		"version":  book.Version,
		"title":    book.Title,
		"writer":   book.Writer,
		"type":     book.Type,
		"update":   book.LastUpdate,
		"chapter":  book.LastChapter,
		"end":      book.EndFlag,
		"read":     book.ReadFlag,
		"download": book.DownloadFlag,
	}
}
