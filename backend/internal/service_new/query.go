package service

import (
	"errors"
	"strconv"

	"github.com/htchan/BookSpider/internal/model"
)

func (serv *ServiceImp) Book(id int, hash string) (*model.Book, error) {
	if hash == "" {
		return serv.rpo.FindBookById(id)
	}

	hashcode, err := strconv.ParseInt(hash, 36, 64)
	if err != nil {
		return nil, errors.New("invalid hash code")
	}

	return serv.rpo.FindBookByIdHash(id, int(hashcode))
}

func (serv *ServiceImp) QueryBooks(title, writer string, limit, offset int) ([]model.Book, error) {
	return serv.rpo.FindBooksByTitleWriter(title, writer, limit, offset)
}

func (serv *ServiceImp) RandomBooks(limit int) ([]model.Book, error) {
	return serv.rpo.FindBooksByRandom(limit)
}
