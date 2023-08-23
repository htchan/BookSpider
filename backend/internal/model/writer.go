package model

import (
	"errors"
	"strings"

	"github.com/rs/zerolog/log"
)

type Writer struct {
	ID   int
	Name string
}

func NewWriter(name string) Writer {
	return Writer{
		ID:   -1,
		Name: name,
	}
}

func (w Writer) Checksum() string {
	if len(w.Name) > 100 {
		log.
			Error().
			Err(errors.New("writer name is too long")).
			Int("writer id", w.ID).
			Str("name", w.Name).
			Msg("generate checksum failed")
		return ""
	}

	writerName := strings.ReplaceAll(w.Name, " ", "")
	return strToShortHex(simplified(writerName))
}
