package flags

import (
	"flag"
)

type Flags struct {
	Operation  *string
	Site       *string
	Id         *int
	MaxThreads *int
}

func NewFlags() *Flags {
	f := new(Flags)
	f.Operation = flag.String("operation", "", "the operation to work on")
	f.Site = flag.String("site", "", "specific site to operate")
	f.Id = flag.Int("id", -1, "specific id to operate")
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
