package service

import (
	"errors"
	"strconv"

	"github.com/htchan/BookSpider/internal/model"
	"github.com/rs/zerolog/log"
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

func (serv *ServiceImp) BookGroup(id int, hash string) (*model.Book, *model.BookGroup, error) {
	var (
		bkIndex  = -1
		group    model.BookGroup
		hashcode int64
		err      error
	)
	if hash == "" {
		group, err = serv.rpo.FindBookGroupByID(id)
		if err != nil {
			return nil, nil, err
		}
	} else {
		hashcode, err = strconv.ParseInt(hash, 36, 64)
		if err != nil {
			return nil, nil, errors.New("invalid hash code")
		}

		group, err = serv.rpo.FindBookGroupByIDHash(id, int(hashcode))
		if err != nil {
			return nil, nil, err
		}
	}

	for i := range group {
		if group[i].Site == serv.name && group[i].ID == id && group[i].HashCode == int(hashcode) {
			bkIndex = i
			break
		}
	}

	groupIDs := make([]string, 0, len(group))
	for _, bk := range group {
		groupIDs = append(groupIDs, bk.String())
	}

	if bkIndex < 0 {
		err := errors.New("books not found")
		log.
			Error().
			Err(err).
			Int("id", id).
			Int64("hashcode", hashcode).
			Strs("book group id", groupIDs).
			Msg("find book group failed")
		return nil, nil, err
	}

	bk := group[bkIndex]
	if bkIndex+1 >= len(group) {
		group = group[:bkIndex]
	} else {
		group = append(group[:bkIndex], group[bkIndex+1:]...)
	}

	return &bk, &group, nil
}

func (serv *ServiceImp) QueryBooks(title, writer string, limit, offset int) ([]model.Book, error) {
	return serv.rpo.FindBooksByTitleWriter(title, writer, limit, offset)
}

func (serv *ServiceImp) RandomBooks(limit int) ([]model.Book, error) {
	return serv.rpo.FindBooksByRandom(limit)
}
