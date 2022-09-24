package arguement

import (
	"errors"
	"flag"
	"strconv"

	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/repo"
	"github.com/htchan/BookSpider/internal/service/site"
)

type Arguement struct {
	Site      *string
	ID        *int
	HashCode  *string
	Operation *string
}

func LoadArgs() *Arguement {
	f := new(Arguement)

	f.Operation = flag.String("operation", "", "operation to work on")
	f.Site = flag.String("site", "", "site name")
	f.ID = flag.Int("id", 0, "book id")
	f.HashCode = flag.String("hash code", "", "book hash code")

	flag.Parse()

	return f
}

func (f *Arguement) IsAllSite() bool {
	return *f.Site == "" && *f.ID == 0 && *f.HashCode == ""
}

func (f *Arguement) IsSite() bool {
	return *f.Site != "" && *f.ID == 0
}

func (f *Arguement) IsBook() bool {
	return *f.Site != "" && *f.ID > 0
}

func (f *Arguement) GetSite(sites map[string]*site.Site) *site.Site {
	return sites[*f.Site]
}

func (f *Arguement) GetBook(sites map[string]*site.Site) *model.Book {
	st := f.GetSite(sites)
	var bk *model.Book
	var err error
	if *f.HashCode != "" {
		bk, err = st.BookFromIDHash(*f.ID, *f.HashCode)
	} else {
		bk, err = st.BookFromID(*f.ID)
	}
	if errors.Is(err, repo.BookNotExist) {
		hash, _ := strconv.ParseInt(*f.HashCode, 36, 64)
		bk = &model.Book{
			Site:     *f.Site,
			ID:       *f.ID,
			HashCode: int(hash),
		}
	} else {
		return nil
	}
	return bk
}
