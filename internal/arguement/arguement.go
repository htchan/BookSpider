package arguement

import (
	"errors"
	"flag"
	"log"
	"strconv"

	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/repo"
	service_new "github.com/htchan/BookSpider/internal/service_new"
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

func (f *Arguement) GetSite(sites map[string]service_new.Service) service_new.Service {
	return sites[*f.Site]
}

func (f *Arguement) GetBook(sites map[string]service_new.Service) *model.Book {
	st := f.GetSite(sites)
	if st == nil {
		return nil
	}

	bk, err := st.Book(*f.ID, *f.HashCode)
	if errors.Is(err, repo.ErrBookNotExist) {
		hash, _ := strconv.ParseInt(*f.HashCode, 36, 64)
		bk = &model.Book{
			Site:     *f.Site,
			ID:       *f.ID,
			HashCode: int(hash),
		}
		log.Printf("[%v] record not found", bk)
	} else if err != nil {
		log.Println(err)
		return nil
	}
	return bk
}
