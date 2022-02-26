package sites

import (
	"github.com/htchan/BookSpider/pkg/flags"
	"errors"
)

func Check(site *Site, args *flags.Flags) (err error) {
	if !args.Valid() || args.IsBook() {
		err = errors.New("invalid arguments")
		return
	}
	if args.IsEverything() || (args.IsSite() && *args.Site == site.Name) {
		return site.database.UpdateBookRecordsStatusByChapter()
	}
	return nil
}
