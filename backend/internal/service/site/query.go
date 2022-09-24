package site

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/htchan/BookSpider/internal/model"
)

func (st *Site) BookFromID(id int) (*model.Book, error) {
	return st.rp.FindBookById(id)
}

func (st *Site) BookFromIDHash(id int, hash string) (*model.Book, error) {
	hashcode, err := strconv.ParseInt(hash, 36, 64)
	if err != nil {
		return nil, errors.New("invalid hash code")
	}
	fmt.Println("hash:", hashcode)
	return st.rp.FindBookByIdHash(id, int(hashcode))
}

func (st *Site) QueryBooks(title, writer string, limit, offset int) ([]model.Book, error) {
	return st.rp.FindBooksByTitleWriter(title, writer, limit, offset)
}

func (st *Site) RandomBooks(limit int) ([]model.Book, error) {
	return st.rp.FindBooksByRandom(limit)
}
