package sites

import (
	"github.com/htchan/BookSpider/pkg/flags"
	"errors"
	// "path/filepath"
	// "time"
	// "os"
)

func (site *Site) validate() (err error) {
	// select some books and try:
	// - fetch the data
	// - download
	// report the success rate or validating regexp work or not
	return
}

func (site *Site) Validate(args flags.Flags) (err error) {
	if !args.Valid() || args.IsBook() {
		err = errors.New("invalid arguments")
		return
	}
	if args.IsEverything() || (args.IsSite() && *args.Site == site.Name) {
		return site.validate()
	}
	return nil
}