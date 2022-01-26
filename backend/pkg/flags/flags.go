package flags

import (
	"flag"
	"strconv"
)

type Flags struct {
	Operation  *string
	Site       *string
	Id         *int
	HashCode   *string
	MaxThreads *int
}

func NewFlags() *Flags {
	f := new(Flags)
	f.Operation = flag.String("operation", "", "the operation to work on")
	f.Site = flag.String("site", "", "specific site to operate")
	f.Id = flag.Int("id", -1, "specific id to operate")
	f.HashCode = flag.String("hash", "", "specific hash code to operate")
	f.MaxThreads = flag.Int("max-threads", -1, "maximum number of threads to carry the process")

	return f
}

func (f *Flags) Load(maxThreadsConfig int) {
	flag.Parse()

	if *f.MaxThreads <= 0 {
		*f.MaxThreads = maxThreadsConfig
	} else if *f.MaxThreads <= 100 {
		*f.MaxThreads = 100
	}
}

func (f *Flags) IsEverything() bool {
	if f.Site == nil { site := "" ; f.Site = &site }
	if f.Id == nil { id := -1 ; f.Id = &id }
	if f.HashCode == nil { hash := "" ; f.HashCode = &hash }
	if f.MaxThreads == nil { maxThreads := -1 ; f.MaxThreads = &maxThreads }
	return *f.Site == "" && *f.Id == -1 &&
		*f.HashCode == "" && *f.MaxThreads == -1
}

func (f *Flags) IsBook() bool {
	if f.Site == nil { site := "" ; f.Site = &site }
	if f.Id == nil { id := -1 ; f.Id = &id }
	if f.HashCode == nil { hash := "" ; f.HashCode = &hash }
	return *f.Site != "" && *f.Id > 0
}

func (f *Flags) GetBookInfo() (string, int, int) {
	if f.Site == nil { site := "" ; f.Site = &site }
	if f.Id == nil { id := -1 ; f.Id = &id }
	if f.HashCode == nil { hash := "" ; f.HashCode = &hash }
	hashCode, err := strconv.ParseInt(*f.HashCode, 36, 64)
	if err != nil {
		hashCode = -1
	}
	return *f.Site, *f.Id, int(hashCode)
}

func (f *Flags) IsSite() bool {
	return f.Site != nil && *f.Site != "" &&
		(f.Id == nil || *f.Id < 0) &&
		(f.HashCode == nil || *f.HashCode == "")
}

func (f *Flags) Valid() bool {
	return f.IsBook() || f.IsSite()
}