package main

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"strings"
	"sync"

	_ "github.com/lib/pq"
	"golang.org/x/sync/semaphore"

	"github.com/caarlos0/env/v6"
	"github.com/siongui/gojianfan"
)

func simplified(s string) string {
	return gojianfan.T2S(s)
}

func strToShortHex(s string) string {
	b := []byte(s)
	encoded := base64.StdEncoding.EncodeToString(b)
	return encoded
}

type DbConf struct {
	Host     string `env:"PSQL_HOST,required" validate:"min=1"`
	Port     string `env:"PSQL_PORT,required" validate:"min=1"`
	User     string `env:"PSQL_USER,required" validate:"min=1"`
	Password string `env:"PSQL_PASSWORD,required" validate:"min=1"`
	Name     string `env:"PSQL_NAME,required" validate:"min=1"`
}

type TitleSource struct {
	Site     string
	ID       int
	Hash     int
	Title    string
	WriterID int
	Name     string
}

type WriterSource struct {
	ID   int
	Name string
}

func main() {
	var conf DbConf
	env.Parse(&conf)
	conn := fmt.Sprintf(
		"host=%v port=%v user=%v password=%v dbname=%v sslmode=disable",
		conf.Host, conf.Port, conf.User, conf.Password, conf.Name,
	)
	fmt.Println(conn)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var titleWg, writerWg sync.WaitGroup

	titleRows, err := db.Query("select books.site, books.id, books.hash_code, books.title from books where status != 'ERROR' and (checksum=$1 or checksum is null)", "")
	if err != nil {
		panic(err)
	}
	defer titleRows.Close()

	sema := semaphore.NewWeighted(100)
	ctx := context.Background()
	for titleRows.Next() {
		titleWg.Add(1)
		sema.Acquire(ctx, 1)
		var src TitleSource
		err := titleRows.Scan(&src.Site, &src.ID, &src.Hash, &src.Title)
		fmt.Println(src.Site, src.ID, src.Hash)
		go func(src TitleSource) {
			defer titleWg.Done()
			defer sema.Release(1)
			cachedSrc := src
			if err != nil {
				log.Println(err)
			}

			if len(src.Title) > 100 {
				log.Printf("[%v-%v-%v] title: %v; writer: %v is too long", cachedSrc.Site, cachedSrc.ID, cachedSrc.Hash, cachedSrc.Title, cachedSrc.Name)
				return
			}

			cachedSrc.Title = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(cachedSrc.Title, " ", ""), "\t", ""), "\n", "")
			titleSum := strToShortHex(simplified(cachedSrc.Title))

			db.Exec("update books set checksum=$1 where site=$2 and id=$3 and hash_code=$4", titleSum, cachedSrc.Site, cachedSrc.ID, cachedSrc.Hash)
		}(src)
	}

	writerRows, err := db.Query("select id, name from writers where checksum=$1 or checksum is null", "")
	if err != nil {
		panic(err)
	}
	defer writerRows.Close()

	for writerRows.Next() {
		writerWg.Add(1)
		sema.Acquire(ctx, 1)
		var src WriterSource
		err := writerRows.Scan(&src.ID, &src.Name)
		fmt.Println(src.ID)
		go func(src WriterSource) {
			defer writerWg.Done()
			defer sema.Release(1)
			cachedSrc := src
			if err != nil {
				log.Println(err)
			}

			if len(src.Name) > 100 {
				log.Printf("[writer] id: %v; ;name: %v is too long", cachedSrc.ID, cachedSrc.Name)
				return
			}

			src.Name = strings.ReplaceAll(cachedSrc.Name, " ", "")
			nameSum := strToShortHex(simplified(fmt.Sprintf("%s", cachedSrc.Name)))

			db.Exec("update writers set checksum=$1 where id=$2", nameSum, cachedSrc.ID)
		}(src)
	}

	titleWg.Wait()
	writerWg.Wait()
}
